// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package mage

import (
	"context"
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"path/filepath"
	"strings"
)

type GoSonarArgs struct {
	Version      string   // Test name used in logging.
	ScannerOpts  []string // Enable race detector.
	Token        string
	Organization string
	ProjectKey   string
	SonarHostUrl string
	QualityGate  bool
	PRBase       string
	PRBranch     string
	PRNumber     string
}

type SonarCloud mg.Namespace

func DefaultGoSonarArgs() GoSonarArgs {
	return GoSonarArgs{
		Version:      SonarVersion,
		ScannerOpts:  strings.Split(SonarScannerOpt, " "),
		Token:        SonarToken,
		Organization: SonarOrg,
		ProjectKey:   SonarProjectKey,
		SonarHostUrl: SonarHostUrl,
		QualityGate:  SonarQualityGate == "true",
		PRBase:       SonarPRBase,
		PRBranch:     SonarPRBranch,
		PRNumber:     SonarPRNumber,
	}
}

// UploadSonarCloud Upload testcoverage to SonarlCloud
func (SonarCloud) UploadSonarCloud() {
	//TODO validate that we have a coverage repot to upload
	//TODO validate that we have the right variables for SonarCloud
	mg.SerialDeps(GoUploadSonarCloud)
}

func GoUploadSonarCloud(ctx context.Context) error {
	params := DefaultGoSonarArgs()
	repoInfo, err := GetProjectRepoInfo()
	if err != nil {
		return fmt.Errorf("failed to determine repo root and package sub dir, %s", err)
	}
	mountPoint := filepath.ToSlash(filepath.Join("/usr", "src"))

	fmt.Println(">> sonarcloud:", params.SonarHostUrl)

	dockerRun := sh.RunCmd("docker", "run")
	var args []string
	sonarImage := fmt.Sprintf("sonarsource/sonar-scanner-cli:%s", params.Version)

	args = append(args,
		"--rm",
		"-v", repoInfo.RootDir+":"+mountPoint,
		"--env", "SONAR_TOKEN="+params.Token,
		sonarImage,

		//Arguments for the CLI
		fmt.Sprintf("-Dsonar.verbose=%t", mg.Verbose()),
		fmt.Sprintf("-Dsonar.qualitygate.wait=%t", params.QualityGate),
	)
	if params.PRBase != "" && params.PRNumber != "" {
		// Parameters that are parsed in order to analyse PR
		fmt.Sprintf(">> running pull request analysis (%s#%s -> %s)\n", params.PRBranch, params.PRNumber, params.PRBase)
		args = append(args,
			fmt.Sprintf("-Dsonar.pullrequest.base=%s", params.PRBase),
			fmt.Sprintf("-Dsonar.pullrequest.branch=%s", params.PRBranch),
			fmt.Sprintf("-Dsonar.pullrequest.key=%s", params.PRNumber),
		)
	} else {
		fmt.Printf(">> Running analysis on the main branch")
	}

	return dockerRun(args...)
}

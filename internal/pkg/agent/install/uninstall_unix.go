// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

//go:build !windows

package install

func isBlockingOnExe(_ error) bool {
	return false
}

func removeBlockingExe(_ error) (string, error) {
	return "", nil
}

// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// CheckIfRoot is a simple function that we can use to determine if
// the ownership of files and directories we create.
var CheckIfRoot = func() bool {
	return os.Getuid() == 0
}

// SudoCallerIds returns the user id and group id of the SUDO caller.
// If either is unset, it returns zero for both values.
// An error is returned if the relevant environment variables
// are not valid integers.
func SudoCallerIds() (uid int, gid int, err error) {
	uidStr := os.Getenv("SUDO_UID")
	gidStr := os.Getenv("SUDO_GID")

	if uidStr == "" || gidStr == "" {
		return 0, 0, nil
	}
	uid, err = strconv.Atoi(uidStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid value %q for SUDO_UID", uidStr)
	}
	gid, err = strconv.Atoi(gidStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid value %q for SUDO_GID", gidStr)
	}
	return
}

// MkdirForUser will call down to os.Mkdir and if the user is running as root,
// the ownership will be changed to the sudo user.  If there is an error
// getting the SudoCallerIds, the directory is removed and an error returned.
func MkdirForUser(dir string, perm os.FileMode) error {
	if err := os.Mkdir(dir, perm); err != nil {
		return err
	}
	if err := ChownToUser(dir); err != nil {
		os.RemoveAll(dir)
		return err
	}
	return nil
}

// MkdirAllForUser will call down to os.MkdirAll and if the user is running as
// root, the ownership will be changed to the sudo user for each directory
// that was created.  If there is an error getting the SudoCallerIds, the
// directory is removed and an error returned.
func MkdirAllForUser(dir string, perm os.FileMode) error {
	// First thing we need to do is to walk the path upwards to find out which
	// directories we are going to be creating, so we can change the ownership
	// of them and remove them on error.
	if IsDirectory(dir) {
		// We are done.
		return nil
	}

	topMostDir := dir
	toCreate := []string{dir}
	for parent := filepath.Dir(dir); !IsDirectory(parent); parent = filepath.Dir(parent) {
		toCreate = append(toCreate, parent)
		topMostDir = parent
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return err
	}
	if err := ChownToUser(toCreate...); err != nil {
		os.RemoveAll(topMostDir)
		return err
	}
	return nil
}

// ChownToUser will attempt to change the ownership of all the paths
// to the user returned by the SudoCallerIds method.  Ownership change
// will only be attempted if we are running as root.
func ChownToUser(paths ...string) error {
	if !CheckIfRoot() {
		return nil
	}
	uid, gid, err := SudoCallerIds()
	if err != nil {
		return err
	}
	for _, path := range paths {
		if err := os.Chown(path, uid, gid); err != nil {
			return err
		}
	}
	return nil
}

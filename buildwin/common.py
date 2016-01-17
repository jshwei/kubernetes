#!/usr/bin/env python

# Copyright 2014 The Kubernetes Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import subprocess
import sys

import utils

KUBE_ROOT=os.path.abspath(os.path.join(os.path.dirname(os.path.realpath(__file__)), os.pardir))
LOCAL_OUTPUT_ROOT=os.path.join(KUBE_ROOT, "_output")
LOCAL_OUTPUT_BINPATH=os.path.join(LOCAL_OUTPUT_ROOT, "bin")

BUILD_TARGETS = {
    "kube-apiserver": os.path.join("cmd, kube-apiserver"), 
    "kube-controller-manager": os.path.join("cmd", "kube-controller-manager"), 
    "kubectl": os.path.join("cmd", "kubectl"),
    "kubelet": os.path.join("cmd", "kubelet"), 
    "kube-proxy": os.path.join("cmd", "kube-proxy"),
    "kube-scheduler": os.path.join("plugin", "cmd", "kube-scheduler")
}


# TODO: ensure python version is 3.5 or abovedir

def kube_build_ensure_go_installed():
    p = subprocess.run(["go", "version"], stdout = subprocess.PIPE, stderr = subprocess.PIPE)
    if p.returncode != 0:
        utils.kube_log_error("Please make sure go tool is installed and is included in the environment variable $PATH")
        sys.exit(p.returncode)


def kube_build_ensure_godeps_installed():
    p = subprocess.run(["godep", "version"], stdout = subprocess.PIPE, stderr = subprocess.PIPE)
    if p.returncode != 0:
        utils.kube_log_error("Please make sure godep tool (https://github.com/tools/godep) is installed and is included in the environment variable $PATH")
        sys.exit(p.returncode)


# TODO
# Since we don't build in container for windows environment,
# and Go depends on folder structure in dependency lookup,
# right now we force the kubernetes to exist under $(GOPATH)/k8s.io/kubernetes.
# The other option is to copy the KUBE_ROOT to k8s.io/kubernetes before build
def kube_build_ensure_correct_directory_setup():
    if not KUBE_ROOT.endswith(os.path.join("k8s.io", "kubernetes")):
        utils.kube_log_error("Please move your kubernetest project from " + KUBE_ROOT + " to $(GOPATH)/k8s.io/kubernetes")
        sys.exit(1)
    
    # bash uses : to separate, and cmd uses ;
    gopath = os.environ.get('GOPATH')
    if not gopath:
        utils.kube_log_error("Please set environment variable $GOPATH to the root of your workspace")
        sys.exit(1)
    
    gopaths = gopath.split(os.pathsep)

    if not any(os.path.join(os.path.abspath(p), "src", "k8s.io", "kubernetes") for p in gopaths):
        utils.kube_log_error("Please make sure your directory from k8s.io/kubernetes is under one of the $GOPATH/src")
        sys.exit(1)


def kube_build_verify_prereqs():
    utils.kube_log_info("Verifying prerequisites...")
    kube_build_ensure_go_installed()
    kube_build_ensure_godeps_installed()
    kube_build_ensure_correct_directory_setup()


def kube_build_target(target):
    utils.kube_log_info("Building target " + target + "...")
    command = ["godep", "go", "build", "-o", os.path.join(LOCAL_OUTPUT_BINPATH, target + ".exe"), '.' + os.sep + BUILD_TARGETS[target]]
    utils.kube_log_info("Building with command " + ' '.join(command))
    os.chdir(KUBE_ROOT)
    p = subprocess.run(command, stdout = subprocess.PIPE, stderr = subprocess.PIPE, universal_newlines = True)
    if p.returncode != 0:
        utils.kube_log_error("Failed while building target " + target)
        utils.kube_log_error("Error output: " + p.stderr)
        sys.exit(p.returncode)
    else:
        utils.kube_log_info("Built target " + target + " successfully")
 
        
def kube_build_all_targets():
    for target in BUILD_TARGETS.keys():
        kube_build_target(target)
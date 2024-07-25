#!/usr/bin/env bash

has_prefix() { case $2 in "$1"*) true;; *) false;; esac; }

function ks() {
	local EXECUTABLE_PATH=""
	local DEFAULT_EXECUTABLE_PATH="./bin/kser"
	declare -a opts

	while test $# -gt 0; do
		opts+=( "$1" )
		shift
	done

	if [ -n "$KSER_EXECUTABLE_PATH" ]; then
		EXECUTABLE_PATH="$KSER_EXECUTABLE_PATH"
	else
		EXECUTABLE_PATH="$DEFAULT_EXECUTABLE_PATH"
	fi
	RESPONSE="$($EXECUTABLE_PATH "${opts[@]}")"
	if [ $? -ne 0 -o -z "$RESPONSE" ]; then
		printf "%s" "$RESPONSE"
		return $?
	fi

	kubeconfig_prefix="__config_"
	namespace_prefix="__namespace_"

	if ! has_prefix "$kubeconfig_prefix" "$RESPONSE" && ! has_prefix "$namespace_prefix" "$RESPONSE"; then
		printf "%s\n" "$RESPONSE"
		return 0
	fi

	if has_prefix "$kubeconfig_prefix" "$RESPONSE"; then
		RESPONSE=${RESPONSE#"$kubeconfig_prefix"}
		remainder="$RESPONSE"
		KUBECONFIG_NAME="${remainder%%,*}"; remainder="${remainder#*,}"
		KUBECONFIG_PATH="${remainder%%,*}"; remainder="${remainder#*,}"

		if [ -z "$KUBECONFIG_NAME" ]; then
			export KS_CURRENT_KUBECONFIG_NAME=""
			export KS_CURRENT_NAMESPACE=""
			export KUBECONFIG=""
			alias k="kubectl"
			printf "Cleared kubeconfig\n"
			return
		fi

		export KS_CURRENT_KUBECONFIG_NAME="$KUBECONFIG_NAME"
		export KUBECONFIG="$KUBECONFIG_PATH"
		alias k="kubectl"
		printf "Switched to kubeconfig %s\n" "$KUBECONFIG_NAME"
		return
	fi

	RESPONSE=${RESPONSE#"$namespace_prefix"}
	namespace="$RESPONSE"

	if [ -z "$namespace" ]; then
		export KS_CURRENT_NAMESPACE=""
		alias k="kubectl"
		printf "Cleared namespace\n"
		return
	fi

	export KS_CURRENT_NAMESPACE="$namespace"
	alias k="kubectl -n $namespace"
	printf "Switched to namespace %s\n" "$namespace"
}

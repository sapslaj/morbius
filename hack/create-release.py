#!/usr/bin/env python3
import argparse
import subprocess
from typing import Iterable, Optional, Union, cast

import semver


def cmd(args: Iterable[str]) -> bytes:
    print("$ ", *args)
    return subprocess.check_output(args)


def get_versions():
    def parse(tag):
        if not tag:
            return None
        try:
            return semver.VersionInfo.parse(tag.removeprefix("v"))
        except ValueError:
            return None

    return filter(lambda v: v, map(parse, cmd(["git", "tag"]).decode().split("\n")))


def get_latest_version(
    versions: Optional[Iterable[semver.VersionInfo]] = None,
) -> semver.VersionInfo:
    if versions is None:
        versions = get_versions()
    sorted_versions = sorted(versions, reverse=True)
    if sorted_versions:
        return sorted_versions[0]
    return semver.VersionInfo(major=0, minor=0, patch=0)


def commit_version_string(commit_message: str) -> bool:
    status = cmd(["git", "status", "main.go"]).decode()
    if "nothing to commit" in status:
        return False
    cmd(["git", "add", "main.go"])
    cmd(["git", "commit", "-m", commit_message])
    return True


def create_git_tag(version: Union[semver.VersionInfo, str], message: Optional[str] = None):
    if isinstance(version, semver.VersionInfo):
        version = str(version)
    if not version.startswith("v"):
        version = f"v{version}"
    if not message:
        message = version
    cmd(["git", "tag", "-a", version, "-m", message])
    return version


def push_git_tag(tag: str, remote: Optional[str] = None):
    if not remote:
        remote = "origin"
    cmd(["git", "push", remote, tag])
    return remote


def bump_version(
    version: Union[semver.VersionInfo, str],
    major: bool = False,
    minor: bool = False,
    patch: bool = False,
    prerelease: bool = False,
    build: bool = False,
) -> semver.VersionInfo:
    if isinstance(version, str):
        version = semver.VersionInfo.parse(version)
    if major:
        version = version.bump_major()
    if minor:
        version = version.bump_minor()
    if patch:
        version = version.bump_patch()
    if prerelease:
        version = version.bump_prerelease()
    if build:
        version = version.bump_build()
    return version


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--from-version",
        default=None,
        help="Specify specific previous version to bump (defaults to latest version in repo)",
    )
    parser.add_argument(
        "--new-version",
        default=None,
        help="Specify new version (only useful for non-supported version bumps)",
    )
    parser.add_argument("--major", action="store_true", help="Major version bump")
    parser.add_argument("--minor", action="store_true", help="Minor version bump")
    parser.add_argument("--patch", action="store_true", help="Patch version bump")
    parser.add_argument("--prerelease", action="store_true", help="Prerelease version bump")
    parser.add_argument("--build", action="store_true", help="Build version bump")
    parser.add_argument("--tag", action="store_true", help="Create git tag")
    parser.add_argument("--tag-message", default=None, help="git tag message (defaults to version)")
    parser.add_argument("--push", action="store_true", help="Push git tag")
    parser.add_argument("--push-remote", default=None, help="git remote to push tag to")
    args = parser.parse_args()

    if args.new_version:
        new_version = semver.VersionInfo.parse(args.new_version)
        print(f"new version given directly")
    else:
        old_version = cast(str, args.from_version) or get_latest_version()
        print(f"old version is {old_version}")
        new_version = bump_version(
            version=old_version,
            major=args.major,
            minor=args.minor,
            patch=args.patch,
            prerelease=args.prerelease,
            build=args.build,
        )
    print(f"new version is {new_version}")
    if args.tag:
        tag = create_git_tag(version=new_version)
        print(f"created git tag {tag}")
        if args.push:
            remote = push_git_tag(tag=tag, remote=args.push_remote)
            print(f"pushed {tag} to {remote}")


if __name__ == "__main__":
    main()

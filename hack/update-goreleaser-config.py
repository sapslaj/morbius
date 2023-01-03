#!/usr/bin/env python3
# Update the generated bits of .goreleaser.yaml
import click
import copy
from ruyaml import YAML

DOCKERS = {
    "common": {
        "build_flag_templates": [
            "--label=org.opencontainers.image.source=https://github.com/sapslaj/{{.ProjectName}}",
            "--label=org.opencontainers.image.created={{.Date}}",
            "--label=org.opencontainers.image.title={{.ProjectName}}",
            "--label=org.opencontainers.image.revision={{.FullCommit}}",
            "--label=org.opencontainers.image.version={{.Version}}",
        ],
        "use": "buildx",
        "extra_files": [
            "minimal-config.yaml",
        ],
    },
    "matrix": [
        {
            "goarch": "amd64",
            "platform": "linux/amd64",
            "arch_tag": "amd64",
        },
        {
            "goarch": "arm64",
            "platform": "linux/arm64/v8",
            "arch_tag": "arm64v8",
        },
        {
            "goarch": "arm",
            "platform": "linux/arm/v7",
            "arch_tag": "armv7",
            "goarm": 7,
        },
        {
            "goarch": "arm",
            "platform": "linux/arm/v6",
            "arch_tag": "armv6",
            "goarm": 6,
        },
    ],
    "manifests": {
        "repository": "ghcr.io/sapslaj/{{ .ProjectName }}",
        "version_tag": "{{ .Version }}",
        "latest_tag": "latest",
    },
}

yaml = YAML()
yaml.indent(mapping=2, sequence=4, offset=2)
yaml.default_flow_style = False


def docker_image_tag(repository, *tag_parts):
    return ":".join(
        [
            repository,
            "-".join(tag_parts),
        ]
    )


def gen_dockers(dockers_config):
    return [
        {
            **copy.deepcopy(dockers_config["common"]),  # Avoid introducing anchors and aliases
            "build_flag_templates": dockers_config["common"]["build_flag_templates"]
            + [
                f"--platform={docker['platform']}",
            ],
            "goos": docker.get("goos", "linux"),
            "goarch": docker["goarch"],
            "goarm": docker.get("goarm", None),
            "image_templates": [
                docker_image_tag(
                    dockers_config["manifests"]["repository"], dockers_config["manifests"][tag], docker["arch_tag"]
                )
                for tag in ["version_tag", "latest_tag"]
            ],
        }
        for docker in dockers_config["matrix"]
    ]


def gen_docker_manifests(dockers_config):
    return [
        {
            "name_template": ":".join([dockers_config["manifests"]["repository"], dockers_config["manifests"][tag]]),
            "image_templates": [
                docker_image_tag(
                    dockers_config["manifests"]["repository"], dockers_config["manifests"][tag], docker["arch_tag"]
                )
                for docker in dockers_config["matrix"]
            ],
        }
        for tag in ["version_tag", "latest_tag"]
    ]


def transform_doc(doc, dockers_config):
    doc["dockers"] = gen_dockers(dockers_config)
    doc["docker_manifests"] = gen_docker_manifests(dockers_config)
    return doc


@click.command()
@click.option("--file", default=".goreleaser.yaml", help="Path for .goreleaser.yaml file")
def main(file):
    with open(file, "r+") as f:
        doc = yaml.load(f)
        doc = transform_doc(doc, DOCKERS)
        f.seek(0)
        f.truncate()
        yaml.dump(doc, f)


if __name__ == "__main__":
    main()

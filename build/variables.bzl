

IMAGE_REPOSITORY_SETTINGS = struct(
    # The container registry for targets which push or pull container images.
    #
    # For example, `gcr.io` for Google Cloud Container Registry or `docker.io`
    # for DockerHub.
    container_registry = "$(container_registry)",

    # Common prefix of container image repositories.
    repository_prefix = "$(image_repo_prefix)",

    # Common tag for container images.
    image_tag = "$(image_tag)",
)


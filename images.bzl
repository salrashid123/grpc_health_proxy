
"""Container image specs.  from https://github.com/world-federation-of-advertisers/cross-media-measurement """


load("//build:variables.bzl", "IMAGE_REPOSITORY_SETTINGS")

_PREFIX = IMAGE_REPOSITORY_SETTINGS.repository_prefix

COMMON_IMAGES = [
    struct(
        name = "default_server_image",
        image = "//:server-image",
        repository = _PREFIX + "/grpc_health_proxy",
    ),
]

ALL_IMAGES = COMMON_IMAGES 
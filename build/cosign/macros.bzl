# from https://github.com/world-federation-of-advertisers/cross-media-measurement """

"""Common macros."""

load("@rules_multirun//:defs.bzl", "command", "multirun")
load("@rules_oci//cosign:defs.bzl", "cosign_sign")

def sign_all(
        name,
        images,
        visibility = None,
        **kwargs):
    """Convenience macro to sign multiple images.

    Args:
      name: a name for the target
      images: dictionary of image repository URL to oci_image target
      visibility: standard visibility attribute
      **kwargs: other args to pass to the resulting target
    """

    for i, (repository_url, image) in enumerate(images.items()):
        sign_name = "{name}_sign_{index}".format(name = name, index = i)
        sign_cmd_name = "{name}_sign_cmd_{index}".format(name = name, index = i)

        cosign_sign(
            name = sign_name,
            image = image,
            repository = repository_url,
            **kwargs
        )
        command(
            name = sign_cmd_name,
            command = ":" + sign_name,
            visibility = ["//visibility:private"],
            **kwargs
        )

    multirun(
        name = name,
        commands = [
            "{name}_sign_cmd_{index}".format(name = name, index = i)
            for i in range(len(images))
        ],
        visibility = visibility,
        **kwargs
    )

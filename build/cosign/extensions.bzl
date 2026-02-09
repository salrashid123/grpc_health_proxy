# from https://github.com/world-federation-of-advertisers/cross-media-measurement """

"extensions for bzlmod"

load("@rules_oci//cosign:repositories.bzl", "cosign_register_toolchains")

def _cosign_impl(module_ctx):
    cosign_register_toolchains(name = "oci_cosign", register = False)
    return module_ctx.extension_metadata()

cosign = module_extension(
    implementation = _cosign_impl,
)

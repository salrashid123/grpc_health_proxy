load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "com_github_alecthomas_kingpin_v2",
        importpath = "github.com/alecthomas/kingpin/v2",
        sum = "h1:f48lwail6p8zpO1bC4TxtqACaGqHYA22qkHjHpqDjYY=",
        version = "v2.4.0",
    )
    go_repository(
        name = "com_github_alecthomas_units",
        importpath = "github.com/alecthomas/units",
        sum = "h1:mimo19zliBX/vSQ6PWWSL9lK8qwHozUj03+zLoEB8O0=",
        version = "v0.0.0-20240927000941-0f3dac36c52b",
    )
    go_repository(
        name = "com_github_beorn7_perks",
        importpath = "github.com/beorn7/perks",
        sum = "h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=",
        version = "v2.3.0",
    )
    go_repository(
        name = "com_github_cncf_xds_go",
        importpath = "github.com/cncf/xds/go",
        sum = "h1:Y8xYupdHxryycyPlc9Y+bSQAYZnetRJ70VMVKm5CKI0=",
        version = "v0.0.0-20251022180443-0feb69152e9f",
    )
    go_repository(
        name = "com_github_creack_pty",
        importpath = "github.com/creack/pty",
        sum = "h1:uDmaGzcdjhF4i/plgjmEsriH11Y0o7RKapEf/LDaM3w=",
        version = "v1.1.9",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane",
        importpath = "github.com/envoyproxy/go-control-plane",
        sum = "h1:K+fnvUM0VZ7ZFJf0n4L/BRlnsb9pL/GuDG6FqaH+PwM=",
        version = "v0.13.5-0.20251024222203-75eaa193e329",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_envoy",
        importpath = "github.com/envoyproxy/go-control-plane/envoy",
        sum = "h1:ixjkELDE+ru6idPxcHLj8LBVc2bFP7iBytj353BoHUo=",
        version = "v1.35.0",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_ratelimit",
        importpath = "github.com/envoyproxy/go-control-plane/ratelimit",
        sum = "h1:/G9QYbddjL25KvtKTv3an9lx6VBE2cnb8wp1vEGNYGI=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_envoyproxy_protoc_gen_validate",
        importpath = "github.com/envoyproxy/protoc-gen-validate",
        sum = "h1:DEo3O99U8j4hBFwbJfrz9VtgcDfUKS7KJ7spH3d86P8=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_go_jose_go_jose_v4",
        importpath = "github.com/go-jose/go-jose/v4",
        sum = "h1:CVLmWDhDVRa6Mi/IgCgaopNosCaHz7zrMeF9MlZRkrs=",
        version = "v4.1.3",
    )
    go_repository(
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:CjnDlHq8ikf6E492q6eKboGOC0T8CDaOvkHCIg8idEI=",
        version = "v1.4.3",
    )
    go_repository(
        name = "com_github_go_logr_stdr",
        importpath = "github.com/go-logr/stdr",
        sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:DrW6hGnjIhtvhOIiAKT6Psh/Kd/ldepEa81DKeiRJ5I=",
        version = "v1.2.5",
    )
    go_repository(
        name = "com_github_golang_jwt_jwt_v5",
        importpath = "github.com/golang-jwt/jwt/v5",
        sum = "h1:pv4AsKCKKZuqlgs5sUmn4x8UlGa0kEVt/puTpKx9vvo=",
        version = "v5.3.0",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=",
        version = "v1.5.4",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        sum = "h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_detectors_gcp",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp",
        sum = "h1:sBEjpZlNHzK1voKq9695PJSX2o5NEXl7/OL3coiIY0c=",
        version = "v1.30.0",
    )
    go_repository(
        name = "com_github_gorilla_mux",
        importpath = "github.com/gorilla/mux",
        sum = "h1:TuBL49tXwgrFYWhqrNgrUNEY92u81SPhu7sTdzQEiWY=",
        version = "v1.8.1",
    )
    go_repository(
        name = "com_github_jpillora_backoff",
        importpath = "github.com/jpillora/backoff",
        sum = "h1:uvFg412JmmHBHw7iwprIxkPMI+sGQ4kzOWsMeHnm2EA=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_json_iterator_go",
        importpath = "github.com/json-iterator/go",
        sum = "h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=",
        version = "v1.1.12",
    )
    go_repository(
        name = "com_github_julienschmidt_httprouter",
        importpath = "github.com/julienschmidt/httprouter",
        sum = "h1:U0609e9tgbseu3rBINet9P48AI/D3oJs4dN7jwJOQ1U=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_klauspost_compress",
        importpath = "github.com/klauspost/compress",
        sum = "h1:c/Cqfb0r+Yi+JtIEq73FWXVkRonBlf0CRNYc8Zttxdo=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_github_kr_pretty",
        importpath = "github.com/kr/pretty",
        sum = "h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_kr_text",
        importpath = "github.com/kr/text",
        sum = "h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_kylelemons_godebug",
        importpath = "github.com/kylelemons/godebug",
        sum = "h1:RPNrshWIDI6G2gRW9EHilWtl7Z6Sb1BR0xunSBf0SNc=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_modern_go_concurrent",
        importpath = "github.com/modern-go/concurrent",
        sum = "h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=",
        version = "v0.0.0-20180306012644-bacd9c7ef1dd",
    )
    go_repository(
        name = "com_github_modern_go_reflect2",
        importpath = "github.com/modern-go/reflect2",
        sum = "h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_munnerz_goautoneg",
        importpath = "github.com/munnerz/goautoneg",
        sum = "h1:C3w9PqII01/Oq1c1nUAm88MOHcQC9l5mIlSMApZMrHA=",
        version = "v0.0.0-20191010083416-a7dc8b61c822",
    )
    go_repository(
        name = "com_github_mwitkow_go_conntrack",
        importpath = "github.com/mwitkow/go-conntrack",
        sum = "h1:KUppIJq7/+SVif2QVs3tOP0zanoHgBEVAwHxUSIzRqU=",
        version = "v0.0.0-20190716064945-2f068394615f",
    )
    go_repository(
        name = "com_github_planetscale_vtprotobuf",
        importpath = "github.com/planetscale/vtprotobuf",
        sum = "h1:GFCKgmp0tecUJ0sJuv4pzYCqS9+RGSn52M3FUwPs+uo=",
        version = "v0.6.1-0.20240319094008-0393e58bdf10",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_prometheus_client_golang",
        importpath = "github.com/prometheus/client_golang",
        sum = "h1:Je96obch5RDVy3FDMndoUsjAhG5Edi49h0RJWRi/o0o=",
        version = "v1.23.2",
    )
    go_repository(
        name = "com_github_prometheus_client_model",
        importpath = "github.com/prometheus/client_model",
        sum = "h1:oBsgwpGs7iVziMvrGhE53c/GrLUsZdHnqNwqPLxwZyk=",
        version = "v0.6.2",
    )
    go_repository(
        name = "com_github_prometheus_common",
        importpath = "github.com/prometheus/common",
        sum = "h1:pIgK94WWlQt1WLwAC5j2ynLaBRDiinoAb86HZHTUGI4=",
        version = "v0.67.5",
    )
    go_repository(
        name = "com_github_prometheus_procfs",
        importpath = "github.com/prometheus/procfs",
        sum = "h1:zUMhqEW66Ex7OXIiDkll3tl9a1ZdilUOd/F6ZXw4Vws=",
        version = "v0.19.2",
    )
    go_repository(
        name = "com_github_rogpeppe_go_internal",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:TMyTOH3F/DB16zRVcYyreMH6GnZZrwQVAoYjRBZyWFQ=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_github_spiffe_go_spiffe_v2",
        importpath = "github.com/spiffe/go-spiffe/v2",
        sum = "h1:l+DolpxNWYgruGQVV0xsfeya3CsC7m8iBzDnMpsbLuo=",
        version = "v2.6.0",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=",
        version = "v1.11.1",
    )
    go_repository(
        name = "com_github_xhit_go_str2duration_v2",
        importpath = "github.com/xhit/go-str2duration/v2",
        sum = "h1:lxklc02Drh6ynqX+DdPyp5pCKLUQpRT8bp8Ydu2Bstc=",
        version = "v2.1.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:pDUj4QMoPejqq20dK0Pg2N4yG9zIkYGdBtwLoEkH9Zs=",
        version = "v0.9.0",
    )
    go_repository(
        name = "dev_cel_expr",
        importpath = "cel.dev/expr",
        sum = "h1:56OvJKSH3hDGL0ml5uSxZmz3/3Pq4tJ+fb1unVLAFcY=",
        version = "v0.24.0",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=",
        version = "v1.0.0-20201130134442-10cb98267c6c",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )
    go_repository(
        name = "in_yaml_go_yaml_v2",
        importpath = "go.yaml.in/yaml/v2",
        sum = "h1:6gvOSjQoTB3vt1l+CU+tSyi/HOjfOjRLJ4YwYZGwRO0=",
        version = "v2.4.3",
    )
    go_repository(
        name = "io_opentelemetry_go_auto_sdk",
        importpath = "go.opentelemetry.io/auto/sdk",
        sum = "h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=",
        version = "v1.2.1",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_detectors_gcp",
        importpath = "go.opentelemetry.io/contrib/detectors/gcp",
        sum = "h1:ZoYbqX7OaA/TAikspPl3ozPI6iY6LiIY9I8cUfm+pJs=",
        version = "v1.38.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel",
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:RkfdswUDRimDg0m2Az18RKOsnI8UDzppJAtj01/Ymk8=",
        version = "v1.38.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:Kl6lzIYGAh5M159u9NgiRkmoMKjvbsKtYRwgfrA6WpA=",
        version = "v1.38.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:l48sr5YbNf2hpCUj/FoGhW9yDkl+Ma+LrVl8qaM5b+E=",
        version = "v1.38.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_metric",
        importpath = "go.opentelemetry.io/otel/sdk/metric",
        sum = "h1:aSH66iL0aZqo//xXzQLYozmWrXxyFkBJ6qT5wthqPoM=",
        version = "v1.38.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:Fxk5bKrDZJUH+AMyyIXGcFAPah0oRcT+LuNtJrmcNLE=",
        version = "v1.38.0",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:+2XxjfsAu6vqFxwGBRcHiMaDCuZiqXGDUDVWVtrFAnE=",
        version = "v0.0.0-20251029180050-ab9386a59fda",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:Jr5R2J6F6qWyzINc+4AM8t5pfUz6beZpHp678GNrMbE=",
        version = "v0.0.0-20260203192932-546029d2fa20",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:K1XZG/yGDJnzMdd/uZHAkVqJE+xIDOcmdSFZkBUicNc=",
        version = "v1.78.0",
    )
    go_repository(
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:fV6ZwhNocDyBLK0dj+fg8ektcVegBBuEolpbTQyBNVE=",
        version = "v1.36.11",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:V6e3FRj+n4dbpw86FJ8Fv7XVOql7TEwpHapKoMJ/GO8=",
        version = "v0.47.0",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:HaW9xtz0+kOcWKwli0ZXy79Ix+UW/vOfmWI5QVd2tgI=",
        version = "v0.31.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:eeHFmOGUTtaaPSGNmjBKpbng9MulQsJURQUAfUwY++o=",
        version = "v0.49.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:hqK/t4AKgbqWkdkcAeI8XLmbK+4m4G5YeQRrmiotGlw=",
        version = "v0.34.0",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:vV+1eWNmZ5geRlYjzm2adRgW2/mcpevXNg50YZtPCE4=",
        version = "v0.19.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:DBZZqJ2Rkml6QMQsZywtnjnnGvHza6BTfYFWY9kjEWQ=",
        version = "v0.40.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:RclSuaJf32jOqZz74CkPA9qFuVTX7vhLlpfj/IGWlqY=",
        version = "v0.39.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:B3njUFyqtHDUI5jMn1YIr5B0IE2U0qck04r6d4KPAxE=",
        version = "v0.33.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:yLkxfA+Qnul4cs9QA3KnlFu0lVmd8JJfoq+E41uSutA=",
        version = "v0.40.0",
    )
    go_repository(
        name = "org_gonum_v1_gonum",
        importpath = "gonum.org/v1/gonum",
        sum = "h1:5+ul4Swaf3ESvrOnidPp4GZbzf0mxVQpDCYUQE7OJfk=",
        version = "v0.16.0",
    )
    go_repository(
        name = "org_uber_go_goleak",
        importpath = "go.uber.org/goleak",
        sum = "h1:2K3zAYmnTNqV73imy9J1T3WC+gmCePx2hEGkimedGto=",
        version = "v1.3.0",
    )

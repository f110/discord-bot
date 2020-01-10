load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_library_dependencies():
    go_repository(
        name = "com_github_bwmarrin_discordgo",
        importpath = "github.com/bwmarrin/discordgo",
        sum = "h1:nA7jiTtqUA9lT93WL2jPjUp8ZTEInRujBdx1C9gkr20=",
        version = "v0.20.2",
    )
    go_repository(
        name = "com_github_gorilla_websocket",
        importpath = "github.com/gorilla/websocket",
        sum = "h1:WDFjx/TMzVgy9VdMMQi2K2Emtwi2QcUQsztZ/zLaH/Q=",
        version = "v1.4.0",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:y6ce7gCWtnH+m3dCjzQ1PCuwl28DDIc3VNnvY29DlIA=",
        version = "v0.0.0-20181030102418-4d3f4d9ffa16",
    )
    go_repository(
        name = "com_github_spf13_pflag",
        importpath = "github.com/spf13/pflag",
        sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
        version = "v1.0.5",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go",
        importpath = "github.com/aws/aws-sdk-go",
        sum = "h1:yr0Lp4bcrIiP8x4JI9wPG+/t4hjdNJghmYJcKX4wh/g=",
        version = "v1.27.2",
    )
    go_repository(
        name = "com_github_jmespath_go_jmespath",
        importpath = "github.com/jmespath/go-jmespath",
        sum = "h1:pmfjZENx5imkbgOkpRUYLnmbU7UEFbjtDA2hxJ1ichM=",
        version = "v0.0.0-20180206201540-c2b33e8439af",
    )

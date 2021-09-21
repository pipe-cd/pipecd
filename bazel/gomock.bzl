# Copyright 2020 The PipeCD Authors.
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

# The MIT License (MIT)
# Copyright © 2018 Jeff Hodges <jeff@somethingsimilar.com>

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the “Software”), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_context")
load("@io_bazel_rules_go//go/private:providers.bzl", "GoLibrary")

_MOCKGEN_TOOL = "@com_github_golang_mock//mockgen"
_MOCKGEN_MODEL_LIB = "@com_github_golang_mock//mockgen/model:go_default_library"

def _gomock_source_impl(ctx):
    args = ["-source", ctx.file.source.path]
    if ctx.attr.package != "":
        args += ["-package", ctx.attr.package]
    args += [",".join(ctx.attr.interfaces)]

    out = ctx.outputs.out
    cmd = ctx.file.mockgen_tool
    go_ctx = go_context(ctx)
    inputs = go_ctx.sdk.headers + go_ctx.sdk.srcs + go_ctx.sdk.tools + [ctx.file.source]

    # We can use the go binary from the stdlib for most of the environment
    # variables, but our GOPATH is specific to the library target we were given.
    ctx.actions.run_shell(
        outputs = [out],
        inputs = inputs,
        tools = [
            cmd,
            go_ctx.go,
        ],
        command = """
           source <($PWD/{godir}/go env) &&
           export PATH=$GOROOT/bin:$PWD/{godir}:$PATH &&
           {cmd} {args} > {out}
        """.format(
            godir = go_ctx.go.path[:-1 - len(go_ctx.go.basename)],
            cmd = "$(pwd)/" + cmd.path,
            args = " ".join(args),
            out = out.path,
        ),
    )

_gomock_source = rule(
    implementation = _gomock_source_impl,
    attrs = {
        "library": attr.label(
            doc = "The target the Go library is at to look for the interfaces in. When this is set and source is not set, mockgen will use its reflect code to generate the mocks. If source is set, its dependencies will be included in the GOPATH that mockgen will be run in.",
            providers = [GoLibrary],
            mandatory = True,
        ),
        "source": attr.label(
            doc = "A Go source file to find all the interfaces to generate mocks for. See also the docs for library.",
            mandatory = False,
            allow_single_file = True,
        ),
        "out": attr.output(
            doc = "The new Go file to emit the generated mocks into",
            mandatory = True,
        ),
        "interfaces": attr.string_list(
            allow_empty = False,
            doc = "The names of the Go interfaces to generate mocks for. If not set, all of the interfaces in the library or source file will have mocks generated for them.",
            mandatory = True,
        ),
        "package": attr.string(
            doc = "The name of the package the generated mocks should be in. If not specified, uses mockgen's default.",
        ),
        "self_package": attr.string(
            doc = "The full package import path for the generated code. The purpose of this flag is to prevent import cycles in the generated code by trying to include its own package. This can happen if the mock's package is set to one of its inputs (usually the main one) and the output is stdio so mockgen cannot detect the final output package. Setting this flag will then tell mockgen which import to exclude.",
        ),
        "mockgen_tool": attr.label(
            doc = "The mockgen tool to run",
            default = Label(_MOCKGEN_TOOL),
            allow_single_file = True,
            executable = True,
            cfg = "host",
            mandatory = False,
        ),
    },
)

def gomock(name, library, out, **kwargs):
    mockgen_tool = _MOCKGEN_TOOL
    if kwargs.get("mockgen_tool", None):
        mockgen_tool = kwargs["mockgen_tool"]

    if kwargs.get("source", None):
        _gomock_source(
            name = name,
            library = library,
            out = out,
            **kwargs
        )
    else:
        _gomock_reflect(
            name = name,
            library = library,
            out = out,
            mockgen_tool = mockgen_tool,
            **kwargs
        )

def _gomock_reflect(name, library, out, mockgen_tool, **kwargs):
    interfaces = kwargs.get("interfaces", None)
    mockgen_model_lib = _MOCKGEN_MODEL_LIB
    if kwargs.get("mockgen_model_library", None):
        mockgen_model_lib = kwargs["mockgen_model_library"]

    prog_src = name + "_gomock_prog"
    prog_src_out = prog_src + ".go"
    _gomock_prog_gen(
        name = prog_src,
        interfaces = interfaces,
        library = library,
        package = kwargs.get("package", None),
        out = prog_src_out,
        mockgen_tool = mockgen_tool,
    )
    prog_bin = name + "_gomock_prog_bin"
    go_binary(
        name = prog_bin,
        srcs = [prog_src_out],
        deps = [library, mockgen_model_lib],
    )
    _gomock_prog_exec(
        name = name,
        interfaces = interfaces,
        library = library,
        package = kwargs.get("package", None),
        out = out,
        prog_bin = prog_bin,
        mockgen_tool = mockgen_tool,
        self_package = kwargs.get("self_package", None),
    )

def _gomock_prog_gen_impl(ctx):
    args = ["-prog_only"]
    if ctx.attr.package != "":
        args += ["-package", ctx.attr.package]
    args += [ctx.attr.library[GoLibrary].importpath]
    args += [",".join(ctx.attr.interfaces)]

    cmd = ctx.file.mockgen_tool
    out = ctx.outputs.out
    ctx.actions.run_shell(
        outputs = [out],
        tools = [cmd],
        command = """
           {cmd} {args} > {out}
        """.format(
            cmd = "$(pwd)/" + cmd.path,
            args = " ".join(args),
            out = out.path,
        ),
    )

_gomock_prog_gen = rule(
    implementation = _gomock_prog_gen_impl,
    attrs = {
        "library": attr.label(
            doc = "The target the Go library is at to look for the interfaces in. When this is set and source is not set, mockgen will use its reflect code to generate the mocks.",
            providers = [GoLibrary],
            mandatory = True,
        ),
        "out": attr.output(
            doc = "The new Go source file put the mock generator code",
            mandatory = True,
        ),
        "interfaces": attr.string_list(
            allow_empty = False,
            doc = "The names of the Go interfaces to generate mocks for. If not set, all of the interfaces in the library or source file will have mocks generated for them.",
            mandatory = True,
        ),
        "package": attr.string(
            doc = "The name of the package the generated mocks should be in. If not specified, uses mockgen's default.",
        ),
        "mockgen_tool": attr.label(
            doc = "The mockgen tool to run",
            default = Label(_MOCKGEN_TOOL),
            allow_single_file = True,
            executable = True,
            cfg = "host",
            mandatory = False,
        ),
    },
)

def _gomock_prog_exec_impl(ctx):
    args = ["-exec_only", ctx.file.prog_bin.path]
    if ctx.attr.package != "":
        args += ["-package", ctx.attr.package]

    if ctx.attr.self_package != "":
        args += ["-self_package", ctx.attr.self_package]

    args += [ctx.attr.library[GoLibrary].importpath]
    args += [",".join(ctx.attr.interfaces)]

    ctx.actions.run_shell(
        outputs = [ctx.outputs.out],
        inputs = [ctx.file.prog_bin],
        tools = [ctx.file.mockgen_tool],
        command = """{cmd} {args} > {out}""".format(
            cmd = "$(pwd)/" + ctx.file.mockgen_tool.path,
            args = " ".join(args),
            out = ctx.outputs.out.path,
        ),
        env = {
            # GOCACHE is required starting in Go 1.12
            "GOCACHE": "./.gocache",
        },
    )

_gomock_prog_exec = rule(
    implementation = _gomock_prog_exec_impl,
    attrs = {
        "library": attr.label(
            doc = "The target the Go library is at to look for the interfaces in. When this is set and source is not set, mockgen will use its reflect code to generate the mocks. If source is set, its dependencies will be included in the GOPATH that mockgen will be run in.",
            providers = [GoLibrary],
            mandatory = True,
        ),
        "out": attr.output(
            doc = "The new Go source file to put the generated mock code",
            mandatory = True,
        ),
        "interfaces": attr.string_list(
            allow_empty = False,
            doc = "The names of the Go interfaces to generate mocks for. If not set, all of the interfaces in the library or source file will have mocks generated for them.",
            mandatory = True,
        ),
        "package": attr.string(
            doc = "The name of the package the generated mocks should be in. If not specified, uses mockgen's default.",
        ),
        "self_package": attr.string(
            doc = "The full package import path for the generated code. The purpose of this flag is to prevent import cycles in the generated code by trying to include its own package. This can happen if the mock's package is set to one of its inputs (usually the main one) and the output is stdio so mockgen cannot detect the final output package. Setting this flag will then tell mockgen which import to exclude.",
        ),
        "prog_bin": attr.label(
            doc = "The program binary generated by mockgen's -prog_only and compiled by bazel.",
            allow_single_file = True,
            executable = True,
            cfg = "host",
            mandatory = True,
        ),
        "mockgen_tool": attr.label(
            doc = "The mockgen tool to run",
            default = Label(_MOCKGEN_TOOL),
            allow_single_file = True,
            executable = True,
            cfg = "host",
            mandatory = False,
        ),
    },
)

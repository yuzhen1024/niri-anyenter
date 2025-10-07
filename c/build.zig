const std = @import("std");


pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const keyevent_monitor_module = b.createModule(.{
        .root_source_file = b.path("keyevent-monitor/main.zig"),
        .target = target,
        .optimize = optimize,
    });
    keyevent_monitor_module.link_libc = true;
    keyevent_monitor_module.linkSystemLibrary("libudev", .{});
    keyevent_monitor_module.linkSystemLibrary("libinput", .{});
    const keyevent_monitor_exe = b.addExecutable(.{
        .name = "keyevent-monitor",
        .root_module = keyevent_monitor_module,
    });

    b.installArtifact(keyevent_monitor_exe);
    // const output = "../bin"
    // const install = b.addInstallArtifact(keyevent_monitor_exe, .{
    //     .dest_dir = .{ 
    //         .override =  .{ 
    //             .custom = output // 无法相对路径
    //         }
    //     }
    // });

    const cmd_run1 = b.addRunArtifact(keyevent_monitor_exe);    // 允许 build run
    if (b.args) |args| {
        cmd_run1.addArgs(args);
    }
    const build_step = b.step("run-1", "run keyevent-monitor");    // 定义run
    build_step.dependOn(&cmd_run1.step);
    build_step.dependOn(b.getInstallStep());
}

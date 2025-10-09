const std = @import("std");


pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});

    const keyevent_monitor_module = b.createModule(.{
        .root_source_file = b.path("keyevent-monitor/main.zig"),
        .target = target,
        .optimize = .ReleaseFast,
    });
    keyevent_monitor_module.strip = true;

    keyevent_monitor_module.link_libc = true;
    keyevent_monitor_module.linkSystemLibrary("libudev", .{});
    keyevent_monitor_module.linkSystemLibrary("libinput", .{});

    const keyevent_monitor_exe = b.addExecutable(.{
        .name = "keyevent-monitor",
        .root_module = keyevent_monitor_module,
    });

    b.installArtifact(keyevent_monitor_exe);
}

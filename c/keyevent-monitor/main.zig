const std = @import("std");
const c = @cImport({
    @cInclude("libinput.h");
    @cInclude("libudev.h");
    @cInclude("poll.h");
    @cInclude("fcntl.h");
    @cInclude("unistd.h");
});

fn open_restricted_fn(path: [*c]const u8, flag: c_int, _: ?*anyopaque) callconv(.c) c_int {
    const fd = c.open(path, flag);
    return fd;
}

fn close_restricted_fn(fd: c_int, _: ?*anyopaque) callconv(.c) void {
    _ = c.close(fd);
}

pub fn main() !void {

    const udev = c.udev_new();
    if (udev == null) {
        @panic("failed to open udev\n");
    }

    const li = c.libinput_udev_create_context( &c.struct_libinput_interface{
        .open_restricted = open_restricted_fn,
        .close_restricted = close_restricted_fn,
    }, null, udev);


    _ = c.libinput_udev_assign_seat(li, "seat0");

    var fds = [_]c.struct_pollfd {
        .{
            .fd = c.libinput_get_fd(li),
            .events = c.POLLIN,
            .revents = 0,
        }
    };

        while (
        // running &&
        c.poll(&fds, 1, -1) > -1
    ) {
        _ = c.libinput_dispatch(li);
        var event: ?*c.libinput_event = c.libinput_get_event(li);
        while (event != null): (event = c.libinput_get_event(li)) {
            // std.debug.print("event-type: {}\n", .{c.libinput_event_get_type(event)});
            if (c.libinput_event_get_type(event) == c.LIBINPUT_EVENT_KEYBOARD_KEY) {
                const keyevent = c.libinput_event_get_keyboard_event(event);
                const keycode = c.libinput_event_keyboard_get_key(keyevent);
                const state = c.libinput_event_keyboard_get_key_state(keyevent);
                // pushKeyEvent(keycode, state);
                try output(keycode, state);
            }

            _ = c.libinput_event_destroy(event);
            _ = c.libinput_dispatch(li);
        }
    }



}

// { key: int, state: bool }
fn output(keycode: u32, state: u32) !void {
    var buf: [64]u8 = undefined;
    const s = try std.fmt.bufPrint(&buf, "{{\"key\": {}, \"state\": {}}}\n", .{keycode, state});
    try std.fs.File.stdout().writeAll(s);
}

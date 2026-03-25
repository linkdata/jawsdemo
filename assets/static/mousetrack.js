var client = {X:0, Y:0, B:0};
var runtime;
onmousemove = function(e) {
    if (typeof jawsVar !== 'undefined') {
        client.X = e.clientX; // or jawsVar('client.X', e.clientX), which would also send the update
        client.Y = e.clientY;
        client.B = e.buttons;
        jawsVar('client');
    }
}

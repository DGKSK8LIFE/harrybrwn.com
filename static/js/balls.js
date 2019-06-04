var canvas = document.getElementById("canvas");
canvas.width = window.innerWidth - 60;
// canvas.width = document.getElementsByClassName("jumbotron").offsetWidth;

// canvas.height = window.innerHeight;
canvas.height = 100;

var ctx = canvas.getContext("2d");

const center = [canvas.width/2, canvas.height/2];
const dxdy = 5;

function drawbox(x, y, r) {
    ctx.rect(x, y, r, r);
    ctx.lineWidth = '1';
    ctx.strokeStyle = "#000000";
    ctx.stroke();
}

class Ball {
    constructor(x, y, r) {
        this.x = x;
        this.y = y;
        this.r = r;
        this.color = 'hsl(' + 360 * Math.random() + ', 50%, 50%)';
        this.theta = 2 * Math.PI;
        this.dy = dxdy * Math.random();
        this.dx = dxdy * Math.random();
    }

    move() {
        var box = this.hitbox();
        if (box.x + box.w > canvas.width || box.x < 0) {
            this.dx = -this.dx;
        }
        if (box.y + box.h > canvas.height || box.y < 0) {
            this.dy = -this.dy;
        }
        this.x += this.dx;
        this.y += this.dy;
    }

    hitbox() {
        var x = this.x - this.r + this.dx;
        var y = this.y - this.r + this.dy;
        var hw = this.r * 2;
        return {x: x, y: y, h: hw, w: hw};
    }

    detectCollision(h2) {
        var h1 = this.hitbox();
        if (h1.x < h2.x + h2.w && h1.x + h1.w > h2.x &&
            h1.y < h2.y + h2.h && h1.y + h1.h > h2.y) {
            return true;
        }
        return false;
    }

    draw() {
        ctx.beginPath();
        ctx.arc(this.x, this.y, this.r, 0, this.theta);
        ctx.fillStyle = this.color;
        ctx.fill();
        ctx.closePath();
    }
}


class ParticleBox {
    constructor(n) {
        this.numBalls = n;
        this.balls = [];
        for (var i = 0; i < this.numBalls; i++) {
            var b = new Ball(
                canvas.width * Math.random(),
                canvas.height * Math.random(), 10
            );
            if (i % 2 == 0) {
                b.dx *= -1;
                b.dy *= -1;
            }
            this.balls.push(b);
        }
    }
    draw() {
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        for (var k = 0; k < this.numBalls; k++) {
            this.balls[k].draw();
            this.balls[k].move();
            console.log(this.balls[k].x, this.balls[k].y);
            // var mouse = {
            //     x: window.event.pageX - 5,
            //     y: window.event.pageY - 5,
            //     w: this.x + 10,
            //     h: this.y + 10
            // }
            // if (this.balls[k].detectCollision(mouse)) {
            //     this.balls[k].x += 1;
            //     this.balls[k].y += 1;
            // }
        }
    }
}

pbox = new ParticleBox(50);
function draw() {
    pbox.draw();
}

setInterval(draw, 10);

kaboom({
	clearColor: [0, 0, 0, 0],
	width: 1000,
	height: 1000,
	scale: 1,
	// fullscreen: true,
	// stretch: true,
	// letterbox: true,
});

loadSprite("player", "sprites/player.png");
loadSprite("wall", "sprites/wall.png");

scene("game", () => {
	const SPEED = 1200;
	const BLOCK_SIZE = 20;
	const MOVE_DELAY = 0.3;
	const dirs = {
		"left": LEFT,
		"right": RIGHT,
		"up": UP,
		"down": DOWN,
	};

	layers(['obj', 'ui'], 'obj')
	camScale(vec2(2))

	function script() {
		return {
			movable: true,
			move_timer: 0,
			update() {
				if (!this.movable) {
					this.move_timer += dt();
					if (this.move_timer > MOVE_DELAY) {
						this.movable = true;
					}
				}
			},
			_move(dir) {
				if (this.movable) {
					player.move(dir.scale(SPEED));
					this.movable = false;
					this.move_timer = 0;
				}
			}
		}
	}

	const maps = [
		[
			'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'x                                                x',
			'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
		]
	]

	const levelCfg = {
		width: BLOCK_SIZE,
		height: BLOCK_SIZE,
		'x': () => [sprite('wall'), scale(0.5), solid(), area()],
	}

	addLevel(maps[0], levelCfg)

	const player = add([
		sprite('player'),
		pos(100, 100),
		solid(),
		area(),
		script(),
	])

	player.action(() => {
		camPos(player.pos)
	})

	// move
	for (const dir in dirs) {
		keyDown(dir, () => {
			player._move(dirs[dir]);
		});
	}
})

go("game");
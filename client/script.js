/*global constrain,Level,noCursor,cursor,ARROW,createButton,loadFont,textFont,buildLevels,Hero,Sprite,Wall,Door,resizeCanvas,world,round,key,abs,rect,ceil,scale,push,pop,frameCount,createCanvas,color,translate,triangle,frameRate,beginShape,endShape,curveVertex,shuffle,sin,cos,floor,rotate,textAlign,LEFT,RIGHT,CENTER,text,textSize,stroke,noStroke,strokeWeight,keyCode,keyIsDown,LEFT_ARROW,RIGHT_ARROW,UP_ARROW,DOWN_ARROW,mouseIsPressed,fill,noFill,mouseX,mouseY,line,ellipse,background,displayWidth,displayHeight,windowWidth,windowHeight,height,width,dist,loadSound,loadImage,image,random,angleMode,RADIANS,DEGREES*/
let world;
let wld = 0;
let lvl = 0;
let gridSize;
let player;
let camera = {
  x: -5,
  y: -5,
  newX: -5,
  newY: -5
};
let players = {};
let images = {};
let hoverHighlightOn = true;
let frameRateSum = 0;
let frameRates = [];
let gameState = "titlescreen";
let host = "https://raw.githubusercontent.com/ohiofi/game-maps/master";
let cameraRatio = 0.03;
let fieldOfView = 11;
let wsUrl = "ws://localhost:8888/world"
let ws;
let x;
let y;
let myFont;
let startButton;
let uid;

const messageType = {
  "JOIN": 0,
  "LEAVE": 1,
  "INIT": 2,
  "MOVE": 3
}

function preload() {

  images.bigtree = loadImage(
    host + "/images/bigtree01.png"
  );
  images.black = loadImage(
    host + "/images/black.png"
  );
  images.brick = loadImage(
    host + "/images/wall2.png"
  );
  images.cat1 = loadImage(
    host + "/images/cat1.png"
  );
  images.cat2 = loadImage(
    host + "/images/cat2.png"
  );
  images.cavebg = loadImage(
    host + "/images/cavetile.png"
  );
  images.chicken1 = loadImage(host + "/images/chicken1.png");
  images.chicken2 = loadImage(host + "/images/chicken2.png");
  images.cliff = loadImage(
    host + "/images/cliff.png"
  );
  images.cobblestone = loadImage(
    host + "/images/cobblestone01.png"
  );
  images.deadtree = loadImage(
    host + "/images/spookytree02.png"
  );
  images.deadtree2 = loadImage(
    host + "/images/spookytree03.png"
  );

  images.gravestone = loadImage(
    host + "/images/pixelgrave02.png"
  );
  images.grassbg = loadImage(
    host + "/images/grass1.png"
  );


  images.hero1 = loadImage(
    host + "/images/heroIdle.png"
  );
  images.hero2 = loadImage(
    host + "/images/heroRun2.png"
  );
  images.hero3 = loadImage(
    host + "/images/heroRun3.png"
  );
  images.hero4 = loadImage(
    host + "/images/heroRun4.png"
  );
  images.rock = loadImage(
    host + "/images/bigrock01.png"
  );

  images.roof = loadImage(
    host + "/images/roof1.png"
  );
  images.sandbg = loadImage(
    host + "/images/sandtile.png"
  );
  images.smalltree = loadImage(
    host + "/images/smalltree2.png"
  );
  images.stump = loadImage(host + "/images/stump2.png");
  images.stump2 = loadImage(host + "/images/stump2.png");
  images.tallgrass = loadImage(
    host + "/images/tallgrass02.png"
  );
  images.water = loadImage(
    host + "/images/water2.png"
  );
  images.water2 = loadImage(
    host + "/images/water2.png"
  );

  myFont = loadFont(
    "https://cdn.glitch.com/f008b3ae-5d6b-4474-9377-414661c88ac7%2FpressStart.ttf?v=1571872754647"
  );
}

function setup() {
  console.log("setup");
  connect()
  createCanvas(windowWidth, windowHeight);
  gridSize = (width + height) * cameraRatio;
}

function draw() {
  if (gameState == "titlescreen") {
    titleScreen();
  }

  if (keyIsDown(UP_ARROW) || keyIsDown(87)) {
    player.newY = round(player.y) - 1;
    player.state = "running";
  } else if (keyIsDown(DOWN_ARROW) || keyIsDown(83)) {
    player.newY = round(player.y) + 1;
    player.state = "running";
  } else if (keyIsDown(LEFT_ARROW) || keyIsDown(65)) {
    player.newX = round(player.x) - 1;
    player.direction = -1;
    player.state = "running";
  } else if (keyIsDown(RIGHT_ARROW) || keyIsDown(68)) {
    player.newX = round(player.x) + 1;
    player.direction = 1;
    player.state = "running";
  }

  if (gameState == "ingame") {
    if (player.x != player.newX || player.y != player.newY) {
      if (player.isRaycastClear()) {
        player.update();
      }

      if (player.msgX != player.newX || player.msgY != player.newY) {
        player.msgX = player.newX;
        player.msgY = player.newY;
        ws.send(JSON.stringify({
          messageType: messageType.MOVE,
          payload: { id: uid, x: player.newX, y: player.newY }
        }));
      }
    } else {
      player.state = "standing";
      player.currentImg = 0;
    }
    drawBackgroundTile();
    world[wld][lvl].showLevelBehindPlayer(player, fieldOfView);
    player.show();
    world[wld][lvl].showLevelInFrontOfPlayer(player, fieldOfView);
    hoverHighlight();
    centerTheCamera();
    showText();
  }
}

function titleScreen() {
  push();
  background("orange");
  stroke(0);
  fill(0);
  textFont(myFont);
  textAlign(CENTER);
  textSize(15);
  text("Click to Start", width * 0.5, height * 0.5);
  pop();
  showText()
  if (keyIsDown(32)) {
    gameState = "ingame";
  }
}

function showText() {
  push()
  strokeWeight(2);
  stroke(255);
  fill(0);
  frameRateSum += frameRate();
  //frameRates.push(frameRateSum += frameRate());
  //if(frameRates.length > 25){
  //  frameRates.splice(0,1)
  //}
  textSize(11);
  textAlign(RIGHT)
  text("FPS: " + round(frameRateSum / frameCount), width - 5, 15);
  textSize(50);
  textAlign(LEFT)
  if (gameState == "ingame") {
    text("HP : " + player.hp, 10, 50);
    text("CNT : " + Object.keys(players).length, 10, 100);
  }
  pop()
}

function resetCamera() {
  camera.x = (player.x * gridSize - gridSize / 2 - width / 2) / gridSize;
  camera.y = ((player.y - .5) * gridSize - gridSize / 2 - height / 2) / gridSize;
  //camera.x = 0.5 - width / 2 / gridSize;
  //camera.y = -height / 2 / gridSize;
  camera.newX = camera.x;
  camera.newY = camera.y;
}

function usingKeys() {
  noCursor();
}

function usingMouse() {
  //cursor(ARROW);
}

function drawBackgroundTile() {
  for (let row = -2 - (camera.y % 1); row < height / gridSize; row++) {
    for (
      let column = -2 - (camera.x % 1);
      column < width / gridSize;
      column++
    ) {
      image(
        world[wld][lvl].backgroundTile,
        column * gridSize,
        row * gridSize,
        gridSize,
        gridSize
      );
    }
  }
}

function hoverHighlight() {
  if (hoverHighlightOn) {
    noStroke();
    fill("rgba(0,0,255,0.2)");
    rect(
      ceil((mouseX + ((camera.x * gridSize) % gridSize)) / gridSize) *
      gridSize -
      ((camera.x * gridSize) % gridSize) -
      gridSize,
      ceil((mouseY + ((camera.y * gridSize) % gridSize)) / gridSize) *
      gridSize -
      ((camera.y * gridSize) % gridSize) -
      gridSize,
      gridSize,
      gridSize
    );
  }
}

function centerTheCamera() {
  camera.newX = (player.x * gridSize - gridSize / 2 - width / 2) / gridSize;
  camera.newY = (player.y * gridSize - gridSize - height / 2) / gridSize;
  camera.x = camera.x + (camera.newX - camera.x) * 0.05;
  camera.y = camera.y + (camera.newY - camera.y) * 0.05;
}

function playerControls() {
  // combo keys
  // if (keyIsDown(UP_ARROW) && keyIsDown(LEFT_ARROW)) {
  //   player.newY = round(player.y) - 0.7;
  //   player.state = "running";
  //   player.newX = round(player.x) - 0.7;
  //   player.direction = -1;
  //   return;
  // }
  // if (keyIsDown(UP_ARROW) && keyIsDown(RIGHT_ARROW)) {
  //   player.newY = round(player.y) - 0.7;
  //   player.state = "running";
  //   player.newX = round(player.x) + 0.7;
  //   player.direction = 1;
  //   return;
  // }
  // if (keyIsDown(DOWN_ARROW) && keyIsDown(LEFT_ARROW)) {
  //   player.newY = round(player.y) + 0.7;
  //   player.state = "running";
  //   player.newX = round(player.x) - 0.7;
  //   player.direction = -1;
  //   return;
  // }
  // if (keyIsDown(DOWN_ARROW) && keyIsDown(RIGHT_ARROW)) {
  //   player.newY = round(player.y) + 0.7;
  //   player.state = "running";
  //   player.newX = round(player.x) + 0.7;
  //   player.direction = 1;
  //   return;
  // }
  // single keys


}

function keyPressed() {
  usingKeys();
  hoverHighlightOn = false;
  if (keyCode == 32) {
    //gameState = "ingame"
    player.toolAnimation();
    // loop thru near objects
    for (let i = ceil(player.y) - 5; i < ceil(player.y) + 5; i++) {
      if (i >= 0 && i < world[wld][lvl].map.length) {
        for (let j = ceil(player.x) - 5; j < ceil(player.x) + 5; j++) {
          if (
            j >= 0 &&
            world[wld][lvl].map[i] &&
            j < world[wld][lvl].map[i].length
          ) {
            if (world[wld][lvl].map[i][j]) {

              //if (world[wld][lvl].objects[i] instanceof TallGrass) {
              if (
                world[wld][lvl].map[i][j].y == player.newY &&
                (world[wld][lvl].map[i][j].x == player.newX ||
                  world[wld][lvl].map[i][j].x == player.newX + player.direction)
              ) {
                // could hurt enemies here eventually?
                //world[wld][lvl].objects.splice(i,1); // if enemy hp <= 0 splice it out
                world[wld][lvl].map[i][j].height = 0; // only affects tallGrass
                world[wld][lvl].map[i][j].growTimer = -900; // only affects tallGrass
              }
              //}
            }
          }
        }
      }
    }
  }
}
function mouseMoved() {
  usingMouse();
}

function mousePressed() {
  // usingMouse();
  // if (gameState == "titlescreen") {
  //   setTimeout(function () {
  //     gameState = "ingame";
  //   }, 50);
  //   return
  // }
  // hoverHighlightOn = true;
  // player.newX = ceil(mouseX / gridSize + camera.x);
  // player.newY = ceil(mouseY / gridSize + camera.y) + 1;
  // if (player.newX > player.x) {
  //   player.direction = 1;
  // } else if (player.newX < player.x) {
  //   player.direction = -1;
  // }
  // player.state = "running";
}

function windowResized() {
  resizeCanvas(windowWidth, windowHeight);
  gridSize = (windowWidth + windowHeight) * 0.05;
}

function connect() {
  ws = new WebSocket(wsUrl);
  ws.onopen = (e) => { console.log("on open: " + e) }
  ws.onclose = (e) => { console.log("on close: " + e) }
  ws.onmessage = (e) => {
    let msg = JSON.parse(e.data);
    let payload = msg.payload;

    switch (msg.messageType) {
      case messageType.INIT:
        // draw player
        uid = payload.player.id;
        world = Level.buildLevels();
        player = new Hero(uid, [images.hero1, images.hero2, images.hero3, images.hero4]);
        player.hp = payload.player.hp
        player.spawn(uid, payload.player.x, payload.player.y);

        // draw other players
        Object.keys(payload.players).map((id) => {
          let player = payload.players[id]
          spawn(id, player.x, player.y)
        })

        resetCamera();
        document.getElementById("p5_loader").style.display = "none";
        break;
      case messageType.MOVE:
        move(payload.id, payload.x, payload.y);
        break;
      case messageType.JOIN:
        if (uid != undefined && payload.id != uid) {
          spawn(payload.id, payload.x, payload.y)
        }
        break;
      case messageType.LEAVE:
        leave(payload.id)
        break;
    }
  }
  ws.onerror = (e) => { console.log("on error: " + e) }
}

function spawn(id, x, y) {
  if (players[id] == undefined) {
    let w = world[0][0];
    let player = new Mover(id, x, y, 1, 1, Sprite.randomDirection(), [images.cat1, images.cat2]);
    players[id] = player;
    w.movers[y].push(player);
  }
}

function move(id, x, y) {
  if (players[id] instanceof Mover) {
    players[id].newX = x;
    players[id].newY = y;
  }
}

function leave(id) {
  if (players[id] instanceof Mover) {
    let w = world[0][0];
    w.movers.forEach(e => {
      e.forEach((mover, i) => {
        if (mover.uid == id) e.splice(i, 1);
      })
    });
  }
}
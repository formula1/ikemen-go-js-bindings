
const UNDEFINED = void(0);

export class GameLoop {
  dead = false;
  running = false;
  milliPerFrame: number

  constructor(
    fps: number,
    private loopFn: ()=>void
  ){
    this.milliPerFrame = 1000 / fps;
  }
    

  play(){
    if(this.dead){
      throw new Error("this loop is dead")
    }
    if(this.running){
      throw new Error("already running")
    }
    this.running = true
  }

  pause(){
    if(this.dead){
      throw new Error("this loop is dead")
    }
    if(!this.running){
      throw new Error("already paused")
    }
    this.running = false
  }

  kill(){
    if(this.dead){
      return;
    }
    this.dead = true;
  }

  async gameLoop(){
    if(this.dead){
      return;
    }
    if(!this.running){
      return;
    }
    while(!this.dead){
      var start = Date.now();
      if(this.running) this.loopFn();
      await new Promise((res)=>{
        const diff = this.milliPerFrame - (Date.now() - start);
        if(diff <= 0) return res(void 0);
        setTimeout(res, diff)
      })
    }
  }
}

function delay(time: number){
  const start = Date.now();
  return new Promise((res)=>{
    const diff = time - (Date.now() - start);
    if(diff <= 0) return res(void 0);
    setTimeout(res, diff)
  })
}

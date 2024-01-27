
class Controller {
  buttons: {
    U: boolean, D: boolean, L: boolean, R: boolean,
    a: boolean, b: boolean, c: boolean,
    x: boolean, y: boolean, z: boolean,
    s: boolean, d: boolean, w: boolean, m: boolean
  } = {
    U: false, D: false, L: false, R: false,
    a: false, b: false, c: false,
    x: false, y: false, z: false,
    s: false, d: false, w: false, m: false
  }
  setButton(key: string, value: boolean){
    if(!(key in this.buttons)) throw new Error("invalid key");
    this.buttons[key] = value;
  }

}
import "wasm/go_wasm_exec.js"
import { setupFileSystem } from "./GameBindings/FileSystem";
import { STATIC_FILES_ORIGIN } from "./constants";

Promise.resolve().then(async ()=>{
  if(!WebAssembly){
    throw new Error("WebAssembly is not available on your browser");
  }
  if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }

  const go = new Go();

  const [goMain, fs, ] = await Promise.all([
    WebAssembly.instantiateStreaming(fetch(STATIC_FILES_ORIGIN + "/main.wasm"), go.importObject),
    setupFileSystem(STATIC_FILES_ORIGIN + "/mugen_base.zip"),
  ]);
  

  console.log(fs.readdir("/"));

  go.run(goMain.instance);
}).catch((e)=>{
  console.error("Failed to build", e);
})

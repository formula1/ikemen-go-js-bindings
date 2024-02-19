import "./wasm/go_wasm_exec.js"
import { setupFileSystem } from "./GameBindings/FileSystem";
import { STATIC_FILES_ORIGIN } from "./constants";

interface IGlobalBinding extends Window {
  IKEMEN_GO_BROWSER_FS: Awaited<ReturnType<typeof setupFileSystem>>
}

Promise.resolve().then(async ()=>{
  console.log("Game Binding Start");
  if(!WebAssembly){
    throw new Error("WebAssembly is not available on your browser");
  }
  if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }

  const GLOBAL_BINDING = window as any as IGlobalBinding;
  const go = new Go();

  const [goMain, fs, ] = await Promise.all([
    WebAssembly.instantiateStreaming(fetch(STATIC_FILES_ORIGIN + "/dist.wasm"), go.importObject),
    setupFileSystem(STATIC_FILES_ORIGIN + "/hidden.mugen_base.zip"),
  ]);
  

  console.log("root dir:", fs.readdirSync("/"));

  GLOBAL_BINDING.IKEMEN_GO_BROWSER_FS = fs;

  console.log("Set fs as -", "IKEMEN_GO_BROWSER_FS", GLOBAL_BINDING.IKEMEN_GO_BROWSER_FS)

  go.run(goMain.instance);
}).catch((e)=>{
  console.error("Failed to build", e);
})

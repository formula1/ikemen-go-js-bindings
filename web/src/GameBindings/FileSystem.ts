console.log("about to import");

import { configure, BFSRequire } from "browserfs";
import { Buffer } from "buffer";
import { promisify } from "util";

console.log("successfuly imported");
export async function setupFileSystem(zipUrl: string){
  const response = await fetch(zipUrl)
  if(!response.ok){
    throw {
      type: "Fetch Error",
      status: response.status,
      text: await response.text()
    }
  }
  const zipArraybuffer = await response.arrayBuffer();
  await promisify(configure)({
    fs: "OverlayFS",
    options: {
      readable: {
        fs: "ZipFS",
        options: {
          zipData: Buffer.from(zipArraybuffer)
        }
      },
      writable: {
        fs: "InMemory"
      }
    }
  })
  return BFSRequire("fs")
}
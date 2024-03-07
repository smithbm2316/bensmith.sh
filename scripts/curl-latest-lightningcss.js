const childProcess = require("child_process");
const { familySync } = require("detect-libc");

// Add `libc` fields only to platforms that have libc(Standard C library).
const triples = [
  {
    name: "x86_64-apple-darwin",
  },
  {
    name: "x86_64-unknown-linux-gnu",
    libc: "glibc",
  },
  {
    name: "x86_64-pc-windows-msvc",
  },
  {
    name: "aarch64-apple-darwin",
  },
  {
    name: "aarch64-unknown-linux-gnu",
    libc: "glibc",
  },
  {
    name: "armv7-unknown-linux-gnueabihf",
  },
  {
    name: "aarch64-unknown-linux-musl",
    libc: "musl",
  },
  {
    name: "x86_64-unknown-linux-musl",
    libc: "musl",
  },
  {
    name: "x86_64-unknown-freebsd",
  },
];
const cpuToNodeArch = {
  x86_64: "x64",
  aarch64: "arm64",
  i686: "ia32",
  armv7: "arm",
};
const sysToNodePlatform = {
  linux: "linux",
  freebsd: "freebsd",
  darwin: "darwin",
  windows: "win32",
};

/** @type {string[]} */
let npmPackages = [];
/** @type {string[]} */
let binaryNames = [];

for (let triple of triples) {
  // Add the libc field to package.json to avoid downloading both
  // `gnu` and `musl` packages in Linux.
  const libc = triple.libc;
  let [cpu, , os, abi] = triple.name.split("-");
  cpu = cpuToNodeArch[cpu] || cpu;
  os = sysToNodePlatform[os] || os;

  let binary = os === "win32" ? "lightningcss.exe" : "lightningcss";
  binaryNames.push(binary);
  let npmPkg = `lightningcss-cli-${os}-${cpu}${abi ? "-" + abi : ""}`;
  npmPackages.push(npmPkg);

  let sysCpu = process.arch;
  let sysOs = process.platform;
  let sysNpmPkg = `lightningcss-cli-${sysOs}-${sysCpu}${abi ? "-" + abi : ""}`;
  let sysLibc = familySync();

  if (npmPkg === sysNpmPkg) {
    if (!triple.libc || triple.libc !== sysLibc) {
      continue;
    }

    fetch("https://registry.npmjs.org/lightningcss-cli").then(
      (res) => res.json(),
    ).then((data) => {
      // let versions = Object.keys(data.versions).reverse();
      let latestVersion = data["dist-tags"].latest;
      console.log(
        `Installing latest version ${latestVersion} of package ${npmPkg}`,
      );

      childProcess.exec(
        `curl -LO https://registry.npmjs.org/${npmPkg}/-/${npmPkg}-${latestVersion}.tgz`,
        (err, stdout, stderr) => {
          if (err) {
            console.error(err);
            return;
          } else if (stderr) {
            console.error(stderr);
            return;
          }

          console.log(stdout);
        },
      );
    });
  }
}

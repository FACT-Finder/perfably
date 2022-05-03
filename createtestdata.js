const fs = require("fs");

const randInt = (min, max) => Math.floor(Math.random() * (max - min + 1) + min);
const randStr = (length) => {
  var result = "";
  var characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  var charactersLength = characters.length;
  for (var i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
};

let semver = [0, 0, 0];

const nextSemver = () => {
  let [major, minor, patch] = semver;

  switch (randInt(0, 20)) {
    case 20:
      major++;
      minor = 0;
      patch = 0;
      break;
    case 19:
    case 18:
    case 17:
      minor++;
      patch = 0;
      break;
    default:
      patch++;
      break;
  }

  semver = [major, minor, patch];
  return semver.join(".");
};

const lines = [];

for (let i = 0; i < 200; i++) {
  const values = {
    "add.avg": randInt(3500, 4500),
    "update.avg": randInt(2300, 2800),
    "delete.avg": randInt(2000, 2500),
    "loc.total": randInt(1000, 1250),
    "loc.covered": randInt(400, 600),
  };
  const commitId = randStr(10);
  const meta = {
    commit: { value: commitId, url: "http://example.org/" + commitId },
    date: {
      value: new Date(
        randInt(Date.now() - 1000 * 60 * 60 * 24 * 365, Date.now())
      )
        .toISOString()
        .slice(0, "2020-02-02".length),
    },
  };

  lines.push(JSON.stringify({ version: nextSemver(), values, meta }));
}

for (let i = 0; i < 200; i++) {
  const values = {
    "add.avg": randInt(4000, 5000),
    "delete.avg": randInt(2000, 2500),
    "update.avg": randInt(1000, 2000),
    "loc.total": randInt(2000, 2200),
    "loc.covered": randInt(1200, 1500),
  };

  lines.push(JSON.stringify({ version: nextSemver(), values }));
}

for (let i = 0; i < 200; i++) {
  const values = {
    "add.avg": randInt(800, 900),
    "delete.avg": randInt(1000, 1400),
    "update.avg": randInt(2000, 2500),
    "loc.total": randInt(3000, 3300),
    "loc.covered": randInt(2800, 3200),
  };

  lines.push(JSON.stringify({ version: nextSemver(), values }));
}

fs.writeFileSync("data/ff.v1.jsonl", lines.join("\n"), { encoding: "utf8" });

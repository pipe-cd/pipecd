const fs = require("fs");
const path = require("path");
const OUTPUT_DIR = process.argv[2];
const files = process.argv.slice(3);

fs.mkdirSync(OUTPUT_DIR);
files.forEach((file) => {
  const basename = path.basename(file);
  const f = fs.readFileSync(file, {
    encoding: "utf-8",
  });

  fs.writeFileSync(
    `${OUTPUT_DIR}/${basename}`,
    f
      .replace(/.*validate_pb.*/g, "")
      .replace(/'.*pkg/g, "'pipe/pkg/app/web")
      .replace(/'.*\/model\//g, "'pipe/pkg/app/web/model/")
  );
});

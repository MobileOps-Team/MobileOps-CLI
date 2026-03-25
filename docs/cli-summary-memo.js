const fs = require("fs");
const {
  Document, Packer, Paragraph, TextRun, Table, TableRow, TableCell,
  Header, Footer, AlignmentType, HeadingLevel, BorderStyle, WidthType,
  ShadingType, PageNumber, LevelFormat
} = require("docx");

const border = { style: BorderStyle.SINGLE, size: 1, color: "CCCCCC" };
const borders = { top: border, bottom: border, left: border, right: border };
const cellMargins = { top: 80, bottom: 80, left: 120, right: 120 };

function headerCell(text, width) {
  return new TableCell({
    borders,
    width: { size: width, type: WidthType.DXA },
    shading: { fill: "1B3A5C", type: ShadingType.CLEAR },
    margins: cellMargins,
    verticalAlign: "center",
    children: [new Paragraph({
      children: [new TextRun({ text, bold: true, color: "FFFFFF", font: "Arial", size: 20 })]
    })]
  });
}

function cell(text, width) {
  return new TableCell({
    borders,
    width: { size: width, type: WidthType.DXA },
    margins: cellMargins,
    children: [new Paragraph({
      children: [new TextRun({ text, font: "Arial", size: 20 })]
    })]
  });
}

function codeCell(text, width) {
  return new TableCell({
    borders,
    width: { size: width, type: WidthType.DXA },
    shading: { fill: "F5F5F5", type: ShadingType.CLEAR },
    margins: cellMargins,
    children: [new Paragraph({
      children: [new TextRun({ text, font: "Courier New", size: 18 })]
    })]
  });
}

const doc = new Document({
  styles: {
    default: {
      document: { run: { font: "Arial", size: 22 } }
    },
    paragraphStyles: [
      {
        id: "Heading1", name: "Heading 1", basedOn: "Normal", next: "Normal", quickFormat: true,
        run: { size: 32, bold: true, font: "Arial", color: "1B3A5C" },
        paragraph: { spacing: { before: 360, after: 200 }, outlineLevel: 0 }
      },
      {
        id: "Heading2", name: "Heading 2", basedOn: "Normal", next: "Normal", quickFormat: true,
        run: { size: 26, bold: true, font: "Arial", color: "2E75B6" },
        paragraph: { spacing: { before: 240, after: 120 }, outlineLevel: 1 }
      },
      {
        id: "Heading3", name: "Heading 3", basedOn: "Normal", next: "Normal", quickFormat: true,
        run: { size: 22, bold: true, font: "Arial", color: "404040" },
        paragraph: { spacing: { before: 200, after: 100 }, outlineLevel: 2 }
      }
    ]
  },
  numbering: {
    config: [
      {
        reference: "bullets",
        levels: [{
          level: 0, format: LevelFormat.BULLET, text: "\u2022", alignment: AlignmentType.LEFT,
          style: { paragraph: { indent: { left: 720, hanging: 360 } } }
        }]
      },
      {
        reference: "numbers",
        levels: [{
          level: 0, format: LevelFormat.DECIMAL, text: "%1.", alignment: AlignmentType.LEFT,
          style: { paragraph: { indent: { left: 720, hanging: 360 } } }
        }]
      },
      {
        reference: "nextSteps",
        levels: [{
          level: 0, format: LevelFormat.DECIMAL, text: "%1.", alignment: AlignmentType.LEFT,
          style: { paragraph: { indent: { left: 720, hanging: 360 } } }
        }]
      }
    ]
  },
  sections: [{
    properties: {
      page: {
        size: { width: 12240, height: 15840 },
        margin: { top: 1440, right: 1440, bottom: 1440, left: 1440 }
      }
    },
    headers: {
      default: new Header({
        children: [new Paragraph({
          border: { bottom: { style: BorderStyle.SINGLE, size: 6, color: "1B3A5C", space: 1 } },
          spacing: { after: 200 },
          children: [
            new TextRun({ text: "MobileOps", bold: true, font: "Arial", size: 18, color: "1B3A5C" }),
            new TextRun({ text: "  |  Internal Technical Memo", font: "Arial", size: 18, color: "808080" })
          ]
        })]
      })
    },
    footers: {
      default: new Footer({
        children: [new Paragraph({
          border: { top: { style: BorderStyle.SINGLE, size: 4, color: "CCCCCC", space: 1 } },
          alignment: AlignmentType.CENTER,
          children: [
            new TextRun({ text: "Confidential  |  Page ", font: "Arial", size: 16, color: "808080" }),
            new TextRun({ children: [PageNumber.CURRENT], font: "Arial", size: 16, color: "808080" })
          ]
        })]
      })
    },
    children: [
      // Title
      new Paragraph({
        spacing: { after: 80 },
        children: [new TextRun({ text: "MobileOps CLI", bold: true, font: "Arial", size: 44, color: "1B3A5C" })]
      }),
      new Paragraph({
        spacing: { after: 200 },
        children: [new TextRun({ text: "Technical Summary & Strategic Rationale", font: "Arial", size: 26, color: "666666" })]
      }),

      // Meta table
      new Table({
        width: { size: 9360, type: WidthType.DXA },
        columnWidths: [2000, 7360],
        rows: [
          new TableRow({ children: [
            new TableCell({ borders, width: { size: 2000, type: WidthType.DXA }, shading: { fill: "F0F4F8", type: ShadingType.CLEAR }, margins: cellMargins,
              children: [new Paragraph({ children: [new TextRun({ text: "Date", bold: true, font: "Arial", size: 20 })] })] }),
            cell("March 24, 2026", 7360)
          ]}),
          new TableRow({ children: [
            new TableCell({ borders, width: { size: 2000, type: WidthType.DXA }, shading: { fill: "F0F4F8", type: ShadingType.CLEAR }, margins: cellMargins,
              children: [new Paragraph({ children: [new TextRun({ text: "Author", bold: true, font: "Arial", size: 20 })] })] }),
            cell("Ricardo", 7360)
          ]}),
          new TableRow({ children: [
            new TableCell({ borders, width: { size: 2000, type: WidthType.DXA }, shading: { fill: "F0F4F8", type: ShadingType.CLEAR }, margins: cellMargins,
              children: [new Paragraph({ children: [new TextRun({ text: "Status", bold: true, font: "Arial", size: 20 })] })] }),
            cell("v0.1.0 - Working prototype, all tests passing", 7360)
          ]}),
          new TableRow({ children: [
            new TableCell({ borders, width: { size: 2000, type: WidthType.DXA }, shading: { fill: "F0F4F8", type: ShadingType.CLEAR }, margins: cellMargins,
              children: [new Paragraph({ children: [new TextRun({ text: "Repo", bold: true, font: "Arial", size: 20 })] })] }),
            cell("MobileOps/mobileops-cli (standalone)", 7360)
          ]})
        ]
      }),

      // Section 1
      new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("1. The Decision: CLI vs MCP Server")] }),

      new Paragraph({
        spacing: { after: 120 },
        children: [new TextRun("We evaluated two approaches for giving technically inclined customers and AI agents programmatic access to MobileOps:")]
      }),

      new Paragraph({
        numbering: { reference: "bullets", level: 0 },
        spacing: { after: 80 },
        children: [
          new TextRun({ text: "MCP Server", bold: true }),
          new TextRun(" \u2014 Purpose-built for AI agent integration (Claude Desktop, etc.). Structured tool definitions, but only useful for AI agents, not human users.")
        ]
      }),
      new Paragraph({
        numbering: { reference: "bullets", level: 0 },
        spacing: { after: 120 },
        children: [
          new TextRun({ text: "CLI (Command-Line Interface)", bold: true }),
          new TextRun(" \u2014 Universal interface for both humans and AI agents. An agent can install it, read --help, and start making API calls. Human power users get the same tool.")
        ]
      }),

      new Paragraph({
        spacing: { after: 120 },
        children: [
          new TextRun({ text: "We chose CLI first", bold: true }),
          new TextRun(", aligning with your vision: customers\u2019 AI agents install the CLI, read the manual, and make API calls. This works today. An MCP wrapper can be added later as a thin layer on top.")
        ]
      }),

      // Section 2
      new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("2. Why Ruby (Not Go)")] }),

      new Paragraph({
        spacing: { after: 120 },
        children: [new TextRun("We initially considered Go for its single-binary distribution advantage (what Basecamp uses). However, Ruby makes more sense for MobileOps:")]
      }),

      new Paragraph({
        numbering: { reference: "bullets", level: 0 },
        spacing: { after: 80 },
        children: [
          new TextRun({ text: "Our team already knows Ruby", bold: true }),
          new TextRun(" \u2014 no context switching or hiring for a new language")
        ]
      }),
      new Paragraph({
        numbering: { reference: "bullets", level: 0 },
        spacing: { after: 80 },
        children: [
          new TextRun({ text: "Faster to build and maintain", bold: true }),
          new TextRun(" \u2014 not learning Go idioms while shipping a product")
        ]
      }),
      new Paragraph({
        numbering: { reference: "bullets", level: 0 },
        spacing: { after: 80 },
        children: [
          new TextRun({ text: "Ruby CLIs are proven", bold: true }),
          new TextRun(" \u2014 Heroku CLI, Rails CLI, and Fastlane are all Ruby")
        ]
      }),
      new Paragraph({
        numbering: { reference: "bullets", level: 0 },
        spacing: { after: 120 },
        children: [
          new TextRun({ text: "Simple distribution", bold: true }),
          new TextRun(" \u2014 an AI agent can run "),
          new TextRun({ text: "gem install mobileops-cli", font: "Courier New", size: 20 }),
          new TextRun(" just as easily as downloading a binary")
        ]
      }),

      // Section 3
      new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("3. What We Built")] }),

      new Paragraph({
        spacing: { after: 120 },
        children: [new TextRun("The CLI is a standalone Ruby gem that wraps our existing REST API. Zero changes to the Rails app were needed.")]
      }),

      new Paragraph({ heading: HeadingLevel.HEADING_2, children: [new TextRun("Tech Stack")] }),

      new Table({
        width: { size: 9360, type: WidthType.DXA },
        columnWidths: [3000, 6360],
        rows: [
          new TableRow({ children: [headerCell("Component", 3000), headerCell("Choice", 6360)] }),
          new TableRow({ children: [cell("CLI Framework", 3000), cell("Thor (same as Rails uses internally)", 6360)] }),
          new TableRow({ children: [cell("HTTP Client", 3000), cell("Faraday (already used in MobileOps-Web)", 6360)] }),
          new TableRow({ children: [cell("Output Formatting", 3000), cell("terminal-table for human display", 6360)] }),
          new TableRow({ children: [cell("Auth", 3000), cell("Existing REST API key pairs (X-Api-Access-Key / X-Api-Secret-Key)", 6360)] }),
          new TableRow({ children: [cell("Testing", 3000), cell("RSpec + WebMock (17 tests, all passing)", 6360)] }),
          new TableRow({ children: [cell("CI/CD", 3000), cell("GitHub Actions: test on Ruby 3.1-3.3, auto-publish to RubyGems on tag", 6360)] }),
        ]
      }),

      new Paragraph({ heading: HeadingLevel.HEADING_2, children: [new TextRun("Commands (v0.1.0)")] }),

      new Table({
        width: { size: 9360, type: WidthType.DXA },
        columnWidths: [4200, 5160],
        rows: [
          new TableRow({ children: [headerCell("Command", 4200), headerCell("Description", 5160)] }),
          new TableRow({ children: [codeCell("mobileops auth login", 4200), cell("Authenticate with API keys", 5160)] }),
          new TableRow({ children: [codeCell("mobileops auth status", 4200), cell("Check connection & credentials", 5160)] }),
          new TableRow({ children: [codeCell("mobileops vessels list", 4200), cell("List all vessels/assets", 5160)] }),
          new TableRow({ children: [codeCell("mobileops vessels get ID", 4200), cell("Get vessel details", 5160)] }),
          new TableRow({ children: [codeCell("mobileops vessels specs ID", 4200), cell("Get vessel specifications", 5160)] }),
          new TableRow({ children: [codeCell("mobileops jobs list", 4200), cell("List jobs (filterable by vessel, date, status)", 5160)] }),
          new TableRow({ children: [codeCell("mobileops jobs get ID", 4200), cell("Get job details", 5160)] }),
          new TableRow({ children: [codeCell("mobileops crew list", 4200), cell("List crew members", 5160)] }),
          new TableRow({ children: [codeCell("mobileops crew get ID", 4200), cell("Get crew member details", 5160)] }),
          new TableRow({ children: [codeCell("mobileops components list", 4200), cell("List components (filterable by vessel)", 5160)] }),
          new TableRow({ children: [codeCell("mobileops work-requests list", 4200), cell("List work requests (filterable by vessel)", 5160)] }),
        ]
      }),

      // Section 4
      new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("4. Agent-First Design")] }),

      new Paragraph({
        spacing: { after: 120 },
        children: [new TextRun("Inspired by the Basecamp CLI, we built with AI agents as a first-class audience:")]
      }),

      new Paragraph({ heading: HeadingLevel.HEADING_3, children: [new TextRun("--json flag on every command")] }),
      new Paragraph({
        spacing: { after: 120 },
        children: [new TextRun("Returns structured JSON with a standard envelope:")]
      }),
      new Paragraph({
        shading: { fill: "F5F5F5", type: ShadingType.CLEAR },
        spacing: { after: 120 },
        children: [new TextRun({ text: '{ "ok": true, "data": {...}, "summary": "Vessel abc123",\n  "breadcrumbs": ["mobileops vessels specs abc123", ...] }', font: "Courier New", size: 18 })]
      }),

      new Paragraph({ heading: HeadingLevel.HEADING_3, children: [new TextRun("Breadcrumbs")] }),
      new Paragraph({
        spacing: { after: 120 },
        children: [new TextRun("Every response suggests what to do next. After getting a vessel, breadcrumbs suggest viewing its specs, components, jobs, or work requests. This lets an AI agent navigate the API without needing documentation.")]
      }),

      new Paragraph({ heading: HeadingLevel.HEADING_3, children: [new TextRun("--help --agent")] }),
      new Paragraph({
        spacing: { after: 120 },
        children: [new TextRun("Returns a complete machine-readable JSON manifest of all commands, parameters, and output shapes. An AI agent runs this once to discover everything the CLI can do.")]
      }),

      // Section 5
      new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("5. Multi-Environment Support")] }),

      new Paragraph({
        spacing: { after: 120 },
        children: [new TextRun("The CLI supports multiple environments with separate credential files:")]
      }),

      new Table({
        width: { size: 9360, type: WidthType.DXA },
        columnWidths: [2400, 4000, 2960],
        rows: [
          new TableRow({ children: [headerCell("Environment", 2400), headerCell("Default Host", 4000), headerCell("Credentials File", 2960)] }),
          new TableRow({ children: [cell("production", 2400), codeCell("https://www.mobileops.at", 4000), codeCell("credentials.json", 2960)] }),
          new TableRow({ children: [cell("staging", 2400), codeCell("https://staging.mobileops.at", 4000), codeCell("credentials.staging.json", 2960)] }),
          new TableRow({ children: [cell("development", 2400), codeCell("http://localhost:3000", 4000), codeCell("credentials.development.json", 2960)] }),
          new TableRow({ children: [cell("test", 2400), codeCell("http://localhost:3001", 4000), codeCell("credentials.test.json", 2960)] }),
        ]
      }),

      new Paragraph({
        spacing: { before: 120, after: 120 },
        children: [new TextRun("Environments can be set via the --env flag, the MOBILEOPS_ENV environment variable, or default to production.")]
      }),

      // Section 6
      new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("6. Architecture")] }),

      new Paragraph({
        spacing: { after: 120 },
        children: [new TextRun("The CLI is a thin HTTP client. It shares no code with the Rails app and requires zero backend changes. It uses the same REST API that third-party integrators use:")]
      }),

      new Table({
        width: { size: 9360, type: WidthType.DXA },
        columnWidths: [3120, 3120, 3120],
        rows: [
          new TableRow({ children: [headerCell("Web App (Browser)", 3120), headerCell("CLI (Terminal)", 3120), headerCell("Chatbot (Future)", 3120)] }),
          new TableRow({ children: [cell("Visual UI for humans", 3120), cell("Text/JSON for power users & agents", 3120), cell("Natural language + deep links", 3120)] }),
        ]
      }),

      new Paragraph({
        spacing: { before: 120, after: 120 },
        alignment: AlignmentType.CENTER,
        children: [new TextRun({ text: "All three interfaces hit the same MobileOps REST API", italics: true, color: "666666" })]
      }),

      // Section 7
      new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("7. Standalone Repo Structure")] }),

      new Paragraph({
        spacing: { after: 120 },
        children: [new TextRun("The CLI lives in its own repo (MobileOps/mobileops-cli) rather than inside MobileOps-Web, because it\u2019s a separate product with its own release cycle, CI pipeline, and distribution:")]
      }),

      new Paragraph({
        shading: { fill: "F5F5F5", type: ShadingType.CLEAR },
        spacing: { after: 120 },
        children: [new TextRun({ text: "mobileops-cli/\n\u251C\u2500\u2500 bin/mobileops              \u2190 Executable\n\u251C\u2500\u2500 lib/mobileops/             \u2190 Core library\n\u2502   \u251C\u2500\u2500 cli.rb, client.rb, config.rb, envelope.rb, formatter.rb\n\u2502   \u2514\u2500\u2500 commands/ (auth, vessels, jobs, crew, components, work_requests)\n\u251C\u2500\u2500 spec/                      \u2190 17 tests (RSpec + WebMock)\n\u251C\u2500\u2500 .github/workflows/         \u2190 CI + auto-release\n\u251C\u2500\u2500 scripts/install.sh         \u2190 curl installer\n\u251C\u2500\u2500 mobileops-cli.gemspec\n\u2514\u2500\u2500 README.md", font: "Courier New", size: 16 })]
      }),

      // Section 8
      new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("8. Next Steps")] }),

      new Paragraph({
        numbering: { reference: "nextSteps", level: 0 },
        spacing: { after: 80 },
        children: [
          new TextRun({ text: "Create GitHub repo", bold: true }),
          new TextRun(" \u2014 Push to MobileOps/mobileops-cli and set up CI")
        ]
      }),
      new Paragraph({
        numbering: { reference: "nextSteps", level: 0 },
        spacing: { after: 80 },
        children: [
          new TextRun({ text: "Publish to RubyGems", bold: true }),
          new TextRun(" \u2014 Add RUBYGEMS_API_KEY secret, tag v0.1.0 to auto-publish")
        ]
      }),
      new Paragraph({
        numbering: { reference: "nextSteps", level: 0 },
        spacing: { after: 80 },
        children: [
          new TextRun({ text: "Add write operations (v0.2.0)", bold: true }),
          new TextRun(" \u2014 Create/update jobs, work requests, crew assignments")
        ]
      }),
      new Paragraph({
        numbering: { reference: "nextSteps", level: 0 },
        spacing: { after: 80 },
        children: [
          new TextRun({ text: "Expand resource coverage", bold: true }),
          new TextRun(" \u2014 Add remaining 30+ REST API resources (deficiencies, inspections, purchase orders, etc.)")
        ]
      }),
      new Paragraph({
        numbering: { reference: "nextSteps", level: 0 },
        spacing: { after: 80 },
        children: [
          new TextRun({ text: "In-app chatbot (future)", bold: true }),
          new TextRun(" \u2014 LLM with function calling to query MobileOps data and return deep links into the web app")
        ]
      }),
      new Paragraph({
        numbering: { reference: "nextSteps", level: 0 },
        spacing: { after: 80 },
        children: [
          new TextRun({ text: "MCP wrapper (future)", bold: true }),
          new TextRun(" \u2014 Thin MCP server on top of the CLI/API for native Claude Desktop integration")
        ]
      }),

      // Footer note
      new Paragraph({ spacing: { before: 400 }, children: [] }),
      new Paragraph({
        border: { top: { style: BorderStyle.SINGLE, size: 4, color: "CCCCCC", space: 1 } },
        spacing: { before: 200 },
        children: [new TextRun({ text: "Questions? Reach out to Ricardo or try it: ", font: "Arial", size: 18, color: "808080" }),
        new TextRun({ text: "gem install mobileops-cli && mobileops --help", font: "Courier New", size: 18, color: "2E75B6" })]
      }),
    ]
  }]
});

Packer.toBuffer(doc).then(buffer => {
  fs.writeFileSync("/Users/ricardo/Code/MobileOps/mobileops-cli/docs/MobileOps-CLI-Summary.docx", buffer);
  console.log("Document created successfully");
});

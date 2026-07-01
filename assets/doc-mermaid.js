(function () {
  function unwrapMermaidBlocks() {
    document.querySelectorAll("pre.mermaid").forEach(function (pre) {
      var code = pre.querySelector("code");
      var div = document.createElement("div");
      div.className = "mermaid";
      div.textContent = code ? code.textContent : pre.textContent;
      pre.replaceWith(div);
    });
  }

  function renderMermaid() {
    if (typeof mermaid === "undefined") {
      return;
    }

    mermaid.initialize({
      startOnLoad: false,
      theme: "base",
      securityLevel: "strict",
      themeVariables: {
        primaryColor: "#eeede8",
        primaryTextColor: "#171a18",
        primaryBorderColor: "#aaa698",
        lineColor: "#5f665f",
        secondaryColor: "#fcfcfa",
        tertiaryColor: "#f5f4f0",
        fontFamily: '"IBM Plex Sans", ui-sans-serif, sans-serif',
      },
      flowchart: {
        htmlLabels: true,
        curve: "basis",
      },
    });

    return mermaid.run({ querySelector: ".mermaid" });
  }

  function init() {
    unwrapMermaidBlocks();
    var blocks = document.querySelectorAll(".mermaid");
    if (blocks.length === 0) {
      return;
    }
    renderMermaid().catch(function (err) {
      console.error("Mermaid render failed:", err);
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();

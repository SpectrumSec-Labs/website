// Progressive enhancement only. The site is fully usable with this file blocked.
// No third-party calls, no inline handlers (CSP: script-src 'self').
(function () {
  "use strict";
  document.documentElement.setAttribute("data-js", "");

  var toggle = document.querySelector(".nav-toggle");
  var nav = document.getElementById("primary-nav");
  if (toggle && nav) {
    var mq = window.matchMedia("(max-width: 820px)");
    var apply = function () {
      if (mq.matches) {
        nav.hidden = toggle.getAttribute("aria-expanded") !== "true";
      } else {
        nav.hidden = false;
      }
    };
    toggle.addEventListener("click", function () {
      var open = toggle.getAttribute("aria-expanded") === "true";
      toggle.setAttribute("aria-expanded", String(!open));
      apply();
    });
    mq.addEventListener("change", function () {
      toggle.setAttribute("aria-expanded", "false");
      apply();
    });
    document.addEventListener("keydown", function (e) {
      if (e.key === "Escape" && toggle.getAttribute("aria-expanded") === "true") {
        toggle.setAttribute("aria-expanded", "false");
        apply();
        toggle.focus();
      }
    });
    apply();
  }

  // Mark external links for assistive tech and safe rel.
  var host = window.location.host;
  document.querySelectorAll('a[href^="http"]').forEach(function (a) {
    if (a.host !== host) {
      a.rel = (a.rel ? a.rel + " " : "") + "noopener";
      if (!a.hasAttribute("aria-label") && a.textContent) {
        a.setAttribute("aria-label", a.textContent.trim() + " (opens in a new tab)");
      }
      a.target = a.target || "_blank";
    }
  });
})();

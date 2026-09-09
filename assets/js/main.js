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

  // Certifications carousel: prev/next buttons just scroll the native,
  // already-scrollable track — swipe/trackpad/keyboard scrolling all work
  // identically with this file blocked; these buttons are pure enhancement
  // (CSS keeps them hidden without [data-js]).
  document.querySelectorAll(".cert-carousel").forEach(function (car) {
    var track = car.querySelector(".cert-carousel__track");
    if (!track) return;
    car.querySelectorAll("[data-cert-dir]").forEach(function (btn) {
      btn.addEventListener("click", function () {
        var dir = parseInt(btn.getAttribute("data-cert-dir"), 10) || 1;
        track.scrollBy({ left: track.clientWidth * 0.8 * dir, behavior: "smooth" });
      });
    });
  });

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

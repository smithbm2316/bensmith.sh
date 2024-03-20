window.customElements.define(
  "return-to-top",
  class ReturnToTopElement extends HTMLElement {
    connectedCallback() {
      try {
        // if the <body> tag is greater than the viewport, we'll execute javascript
        // below to unhide the <return-to-top> button which is hidden by default.
        // we should only show the tag if the contents of the page exceed the
        // viewport's height. it stays disabled inside of the default HTML
        if (document.querySelector("body").offsetHeight <= window.innerHeight) {
          return;
        }

        let button = this.querySelector(":scope button");
        if (!button) {
          return;
        }
        button.classList.remove("hidden");
        button.ariaHidden = null;

        // set up a click handler for our button to scroll to the top of the page
        button.addEventListener("click", () => {
          window.scrollTo({
            // make sure to be inclusive to users who prefer reduced motion
            behavior:
              window.matchMedia("(prefers-reduced-motion: no-preference)")
                ? "smooth"
                : "auto",
            top: 0,
          });
        });
      } catch (err) {
        if (this.dataset.debug === "true") {
          console.log(err);
        }
      }
    }
  },
);

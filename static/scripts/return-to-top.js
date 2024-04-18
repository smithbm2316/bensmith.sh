window.customElements.define(
  "return-to-top",
  class ReturnToTopElement extends HTMLElement {
    connectedCallback() {
      try {
        let button = this.querySelector(":scope button");
        if (!button) {
          return;
        }

        let returnToTop = {
          intersectionElement: document.querySelector(
            "#return-to-top-intersection",
          ),
          isVisible: false,
        };
        if (!returnToTop.intersectionElement) {
          return;
        }

        // setup an intersection observer that will watch the element with the
        // id of `#return-to-top-intersection` to show/hide our <return-to-top>
        // button as we scroll up/down the page
        let observer = new IntersectionObserver((entries) => {
          for (let entry of entries) {
            if (!returnToTop.isVisible && entry.boundingClientRect.y < 0) {
              button.classList.remove("hidden");
              button.ariaHidden = null;
              returnToTop.isVisible = true;
            } else if (
              returnToTop.isVisible && entry.boundingClientRect.y >= 0
            ) {
              button.classList.add("hidden");
              button.ariaHidden = "true";
              returnToTop.isVisible = false;
            }
          }
        }, {
          rootMargin: "0px",
          threshold: 1.0,
        });
        observer.observe(returnToTop.intersectionElement);

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

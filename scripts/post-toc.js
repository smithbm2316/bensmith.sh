window.customElements.define(
  "post-toc",
  class PostTOCElement extends HTMLElement {
    /**
     * If our table of contents is currently active in its own column due to
     * our `article-container` container query, then set the `<details>` element
     * to be `open="true"` by default
     */
    autoExpandOnDesktop() {
      try {
        let detailsElement = this.querySelector(":scope details");
        if (
          getComputedStyle(detailsElement).getPropertyValue(
            "--auto-open-toc",
          ) === "true"
        ) {
          detailsElement.open = true;
        }
      } catch (err) {
        if (this.dataset.debug === "true") {
          console.log(err);
        }
      }
    }

    connectedCallback() {
      this.autoExpandOnDesktop();
    }
  },
);

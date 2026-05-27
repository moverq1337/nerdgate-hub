function syncTargetPanels() {
  const selected = document.querySelector('input[name="target_mode"]:checked');
  const mode = selected ? selected.value : "manual";

  document.querySelectorAll("[data-target-panel]").forEach((panel) => {
    const active = panel.dataset.targetPanel === mode;
    panel.hidden = !active;
    panel.querySelectorAll("input, select").forEach((control) => {
      control.disabled = !active || control.dataset.unavailable === "true";
      if (control.name === "target_url" || control.name === "docker_target") {
        control.required = active && control.dataset.unavailable !== "true";
      }
    });
  });
}

function setupTabs() {
  document.querySelectorAll("[data-tabs]").forEach((container) => {
    const tabs = container.querySelectorAll(".panel-tab");
    const panels = container.querySelectorAll(".panel-tabpanel");
    tabs.forEach((tab) => {
      tab.addEventListener("click", () => {
        const key = tab.dataset.tab;
        tabs.forEach((t) => {
          const active = t.dataset.tab === key;
          t.classList.toggle("is-active", active);
          t.setAttribute("aria-selected", active ? "true" : "false");
        });
        panels.forEach((p) => {
          const active = p.dataset.tabpanel === key;
          p.classList.toggle("is-active", active);
          if (active) {
            p.removeAttribute("hidden");
          } else {
            p.setAttribute("hidden", "");
          }
        });
      });
    });
  });
}

function setupFileDrop() {
  document.querySelectorAll(".file-drop").forEach((drop) => {
    const input = drop.querySelector('input[type="file"]');
    const label = drop.querySelector(".file-drop-filename");
    if (!input || !label) {
      return;
    }
    const placeholder = drop.dataset.emptyText || "No file selected";

    const update = () => {
      if (input.files && input.files.length > 0) {
        label.textContent = input.files[0].name;
        drop.classList.add("has-file");
      } else {
        label.textContent = placeholder;
        drop.classList.remove("has-file");
      }
    };

    input.addEventListener("change", update);

    ["dragenter", "dragover"].forEach((ev) => {
      drop.addEventListener(ev, (e) => {
        e.preventDefault();
        drop.classList.add("is-drag");
      });
    });

    ["dragleave", "dragend"].forEach((ev) => {
      drop.addEventListener(ev, () => {
        drop.classList.remove("is-drag");
      });
    });

    drop.addEventListener("drop", (e) => {
      e.preventDefault();
      drop.classList.remove("is-drag");
      if (!e.dataTransfer || !e.dataTransfer.files || !e.dataTransfer.files.length) {
        return;
      }
      try {
        const dt = new DataTransfer();
        dt.items.add(e.dataTransfer.files[0]);
        input.files = dt.files;
      } catch {
        /* DataTransfer constructor unavailable in some environments */
      }
      update();
    });

    update();
  });
}

document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll('input[name="target_mode"]').forEach((input) => {
    input.addEventListener("change", syncTargetPanels);
  });
  syncTargetPanels();

  setupTabs();
  setupFileDrop();

  document.querySelectorAll("form").forEach((form) => {
    form.addEventListener("submit", (event) => {
      const message = form.dataset.confirm;
      if (message && !window.confirm(message)) {
        event.preventDefault();
        return;
      }

      const submitter = event.submitter || form.querySelector('button[type="submit"]');
      if (!submitter) {
        return;
      }

      const loadingLabel = submitter.dataset.loadingLabel || "Working...";
      submitter.dataset.originalLabel = submitter.textContent;
      submitter.textContent = loadingLabel;
      submitter.setAttribute("aria-busy", "true");
      submitter.disabled = true;
    });
  });
});

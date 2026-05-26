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

document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll('input[name="target_mode"]').forEach((input) => {
    input.addEventListener("change", syncTargetPanels);
  });
  syncTargetPanels();
});


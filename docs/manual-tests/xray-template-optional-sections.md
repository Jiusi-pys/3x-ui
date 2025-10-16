# Xray Template Optional Sections Regression

These steps verify that editing and saving templates lacking optional sections
(outbounds, routing rules, and `policy.system`) continues to work after the
safety guards introduced in this change.

1. Open the panel and navigate to **Settings → Xray**.
2. Click **Edit JSON** (Advanced Template tab) and remove the `routing.rules`,
   `outbounds`, and `policy.system` sections from the JSON payload. Leave other
   configuration intact and save the changes.
3. Press **Save** and confirm the request succeeds (no toast errors) and the
   editor automatically re-populates default empty objects/arrays for the
   removed sections.
4. Switch between the **Routing** and **Outbounds** tabs to ensure their tables
   render without JavaScript errors and display empty states instead of
   crashing.
5. Reload the page to make sure the saved configuration remains valid and the
   UI still loads without errors.


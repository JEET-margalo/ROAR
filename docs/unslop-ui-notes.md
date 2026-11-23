# ROAR UI polish notes

ROAR's dashboard was reviewed using the `unslop` skill from `https://github.com/heydamir/unslop`.
The skill is licensed CC BY 4.0; this adaptation credits Damir Sokolovsky.

## Direction

- Product: a scoped security testing dashboard for authorized assets
- Character: calm, exact, operational
- World: an instrument panel used during a security review
- Moments: live scan state and evidence-backed findings
- Signature: a restrained status indicator that explains what the scanner actually did

## Applied polish

- Flattened the logo tile and wordmark; removed gradient branding.
- Removed glow from the primary action and logo.
- Kept one cyan action accent while preserving semantic colors for findings.
- Reduced unnecessary uppercase/letter-spaced chrome labels.
- Kept cards quiet with a small neutral shadow rather than decorative glow.
- Added `prefers-reduced-motion` support.
- Preserved existing routes, features, data, and dark theme.

## Deliberately unchanged

- Scanner status colors remain semantic.
- The coverage snapshot remains visually distinct because it communicates scan evidence.
- No decorative imagery or new claims were added.
- No workflow or security control was removed.

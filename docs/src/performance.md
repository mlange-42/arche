# Performance tips

## Random/World access

Prefer accessing few from many over many from view:
- Query over children with world access to parents is better than query over parents.
- Avoids random access over a large memory section.
- May make code more complicated and/or require additional components for temporary data.
- Consider using Arche's entity relations feature!

## Entity relations

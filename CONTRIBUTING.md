# Contributing Guidelines

This is a starter kit repository. Contributions should maintain the clean structure and best practices.

## Before Starting

1. Read the documentation:
   - `BACKEND_API_GUIDE.md` - API architecture & patterns
   - `STARTER_KIT_SETUP.md` - Development workflow
   - `STARTER_KIT_SUMMARY.md` - Quick reference

2. Follow the established patterns in template files:
   - `app/Http/Controllers/CONTROLLER_TEMPLATE.go`
   - `app/Models/MODEL_TEMPLATE.go`
   - `app/Http/Requests/EXAMPLE_REQUEST_TEMPLATE.go`
   - `database/migrations/MIGRATION_TEMPLATE.sql`

## Code Style

- Use `gofmt` for code formatting
- Follow Go conventions and idioms
- Use meaningful variable and function names
- Add comments for exported functions and complex logic

## File Organization

- Controllers go in `app/Http/Controllers/`
- Models go in `app/Models/`
- Request validation go in `app/Http/Requests/`
- Migrations go in `database/migrations/`
- Seeders go in `database/seeders/`

## Commit Messages

Use clear, descriptive commit messages:

```
type(scope): subject

body

footer
```

Types:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Code style changes (formatting, etc.)
- `refactor` - Code refactoring
- `test` - Adding tests
- `chore` - Maintenance tasks

Example:
```
feat(auth): add token refresh endpoint

Add refresh token endpoint to refresh expired access tokens.
Update JWT utility to support token refresh flow.

Closes #123
```

## Testing

Before submitting:
1. Test your changes locally
2. Run `go mod tidy` to update dependencies
3. Verify all endpoints work as expected
4. Check for any lint issues

## Documentation

When adding new features:
1. Update relevant documentation files
2. Add code comments where needed
3. Update `BACKEND_API_GUIDE.md` if API changes
4. Update `README.md` if user-facing changes

## Git Workflow

1. Create a feature branch
2. Make your changes
3. Commit with clear messages
4. Push to your fork
5. Create a pull request

## Cleanup

To reset to starter kit state:
```bash
./clean_starter.sh
```

This removes all implementations while keeping infrastructure.

---

Thank you for contributing to maintaining a clean, professional starter kit!

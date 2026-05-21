# Seosementara

## Cursor Cloud specific instructions

### Overview

This is a WordPress plugin/theme project ("Seosementara" - Mass Onboarding SEO) designed for 50–55 domains on shared hosting, where each domain can have thousands of posts/pages. The repository is currently in initial state.

### Development Environment

- **WordPress 7.0** installed at `/var/www/wordpress`
- **PHP 8.3.6** with extensions: mysql, xml, mbstring, curl, zip, gd, intl, bcmath
- **MariaDB 10.11** with database `wordpress` (root user, no password, localhost)
- **WP-CLI 2.12** for WordPress management
- **PHPCS 3.13** with WordPress Coding Standards (WordPress, WordPress-Core, WordPress-Docs, WordPress-Extra)
- **PHPUnit 12.5** for automated testing
- **Composer 2.9** for PHP dependency management

### Starting Services

Before development, start these services:

```bash
sudo service mariadb start
php -S localhost:8080 -t /var/www/wordpress &
```

### Plugin Development

The workspace (`/workspace`) is symlinked to `/var/www/wordpress/wp-content/plugins/seosementara`. Any PHP files added to the workspace will be available as a WordPress plugin.

### Linting

```bash
phpcs --standard=WordPress --extensions=php /workspace/
```

### Testing

```bash
phpunit
```

For WordPress-specific integration tests, use `wp scaffold plugin-tests` to set up the test suite.

### WP-CLI Usage

Always use `--allow-root` flag when running WP-CLI commands:

```bash
cd /var/www/wordpress && wp <command> --allow-root
```

### Key Constraints (from project rules)

- Minimize database queries (50-55 domains on shared hosting)
- Use `wpmo_extend_execution_time()` before heavy operations
- Batch operations, avoid N+1 queries
- Cache with transients where appropriate
- Never load all posts without pagination (`posts_per_page` limit)

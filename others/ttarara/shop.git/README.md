# Shop

Rails e-commerce starter where users sign up (Devise), list products (“ads”), and use a shopping cart. This repo extends the original exercise with registration strong params, seller/ownership helpers, cart persistence for guests and signed-in users, and updated styling.

## Stack

- Ruby **3.0.0**
- Rails **6.1.x** (project targets 6.1.3)
- SQLite
- [Devise](https://github.com/heartcombo/devise) (authentication)
- [Bulma](https://bulma.io/) (via `bulma-rails`) + SCSS
- [CarrierWave](https://github.com/carrierwaveuploader/carrierwave) (product images)

## Features

- **User registration** with Devise and a custom `RegistrationsController` permitting `name`, `email`, `password`, and `password_confirmation` on sign up; account updates also permit `current_password`.
- **Products** with brand, finish, condition, title, price, model, description, and image upload.
- **Seller display** and **edit/delete only for the owner** (helper + controller guard).
- **Shopping cart**: add items (including from the product index), per-line remove, total, empty cart (redirects home), flash messages for add/remove, navbar count.
- **`CurrentCart` concern**: guest cart in session is merged into the user’s cart after sign in.

## Prerequisites

- Ruby 3.0.0 (e.g. [rbenv](https://github.com/rbenv/rbenv))
- Bundler 2.x
- Node + Yarn if you use Webpacker fully (this app loads `javascript_pack_tag`; a minimal pack exists for UJS/Turbolinks)

On Ubuntu, typical build packages:

```bash
sudo apt install -y build-essential libssl-dev zlib1g-dev libreadline-dev \
  libyaml-dev libsqlite3-dev sqlite3 libffi-dev libgdbm-dev git curl
```

## Setup

From the application directory (the folder that contains `Gemfile`):

```bash
cd "/path/to/shop/shop"   # folder that contains Gemfile (e.g. ~/Documents/shop/shop)

bundle install
bundle exec rails db:create db:migrate db:seed
bundle exec rails server
```

Open [http://localhost:3000](http://localhost:3000).

### Performance “ms” panel (rack-mini-profiler)

In development, the **rack-mini-profiler** gem can show a small timing sidebar. It is **not** required for the app; it is for debugging speed.

By default it is **off** (see `config/initializers/rack_mini_profiler.rb`). To turn it on temporarily:

```bash
RACK_MINI_PROFILER=1 bundle exec rails server
```

### If `bundle install` fails on `mimemagic`

Older `Gemfile.lock` entries may pin `mimemagic` versions that were yanked from RubyGems. Try:

```bash
bundle update mimemagic
bundle install
```

## Configuration

- **Devise**: `config/initializers/devise.rb`, routes in `config/routes.rb` (`devise_for :users, controllers: { registrations: 'registrations' }`).
- **Mailer URLs** (for password reset, etc.): set `config.action_mailer.default_url_options` in `config/environments/development.rb` if you enable mail delivery.

## Project layout (high level)

| Area | Path |
|------|------|
| Custom registration params | `app/controllers/registrations_controller.rb` |
| Products + ownership | `app/controllers/products_controller.rb`, `app/helpers/products_helper.rb` |
| Cart | `app/models/cart.rb`, `app/models/cart_item.rb`, `app/controllers/carts_controller.rb`, `app/controllers/cart_items_controller.rb` |
| Guest → user cart | `app/models/concerns/current_cart.rb`, included in `app/controllers/application_controller.rb` |
| Routes | `config/routes.rb` |
| Main layout & nav | `app/views/layouts/application.html.erb` |
| Global + product UI styles | `app/assets/stylesheets/application.scss`, `app/assets/stylesheets/products.scss` |

## Styling notes

- Global theme and navbar tweaks live in **`app/assets/stylesheets/application.scss`** (after `@import "bulma"`).
- The **“Shop”** word in the top bar is styled with **`.navbar-brand .title.is-5`** in that file—change `color`, `font-size`, or `background` there if you only want to adjust the logo.

## Git remote (example)

```bash
git remote add origin https://platform.zone01.gr/git/ttarara/shop.git
git branch -M main   # if your default branch should be main
git push -u origin main
```

Use your real remote URL and branch name.

## License

Educational / assignment use unless otherwise specified by your course or employer.

class Cart < ApplicationRecord
  belongs_to :user, optional: true
  has_many :cart_items, dependent: :destroy

  def add_product(product)
    item = cart_items.find_or_initialize_by(product: product)
    item.quantity = item.new_record? ? 1 : item.quantity + 1
    item.save
    item
  end

  def total
    cart_items.includes(:product).sum { |item| item.subtotal }
  end

  def items_count
    cart_items.sum(:quantity)
  end

  def merge_from!(other_cart)
    return if other_cart.nil? || other_cart == self

    other_cart.cart_items.find_each do |item|
      existing = cart_items.find_or_initialize_by(product: item.product)
      existing.quantity = existing.new_record? ? item.quantity : existing.quantity + item.quantity
      existing.save!
    end

    other_cart.destroy
  end
end

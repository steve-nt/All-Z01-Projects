module ProductsHelper
  def product_author(product)
    product.user&.name.presence || "Unknown seller"
  end

  def can_manage_product?(product)
    user_signed_in? && product.user == current_user
  end
end

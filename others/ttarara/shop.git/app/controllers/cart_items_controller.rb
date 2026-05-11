class CartItemsController < ApplicationController
  def create
    product = Product.find(params[:product_id])
    @cart.add_product(product)
    redirect_back fallback_location: products_path, notice: "Added to your cart"
  end

  def destroy
    item = @cart.cart_items.find(params[:id])
    item.destroy
    redirect_to cart_path, notice: "Removed from your cart"
  end
end

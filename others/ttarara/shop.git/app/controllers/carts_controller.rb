class CartsController < ApplicationController
  def show; end

  def empty
    @cart.cart_items.destroy_all
    redirect_to root_path, notice: "Cart emptied."
  end
end

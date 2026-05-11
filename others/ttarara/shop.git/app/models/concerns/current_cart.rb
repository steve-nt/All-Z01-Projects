module CurrentCart
  extend ActiveSupport::Concern

  private

  def set_cart
    @cart = load_cart
  end

  def load_cart
    session_cart = Cart.find_by(id: session[:cart_id])

    if user_signed_in?
      user_cart = current_user.cart || current_user.create_cart
      if session_cart && session_cart != user_cart
        user_cart.merge_from!(session_cart)
      end

      session[:cart_id] = user_cart.id
      user_cart
    else
      session_cart || create_guest_cart
    end
  end

  def create_guest_cart
    cart = Cart.create
    session[:cart_id] = cart.id
    cart
  end
end

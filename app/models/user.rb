# frozen_string_literal: true

# User model
class User < ApplicationRecord
  CONFIRMATION_TOKEN_EXPIRATION = 10.minutes
  PASSWORD_RESET_TOKEN_EXPIRATION = 10.minutes
  MAILER_FROM_EMAIL = 'support@httpmock.dev'
  EMAIL_PASSWORD_PROVIDER = 'email/password'

  has_secure_password

  has_many :active_sessions, dependent: :destroy
  has_many :projects, dependent: :destroy
  has_many :endpoints, dependent: :destroy

  attr_accessor :current_password

  before_save :downcase_email
  before_save :downcase_unconfirmed_email
  before_create :generate_api_key
  before_create :set_provider

  validates :email, presence: true, format: { with: URI::MailTo::EMAIL_REGEXP }, uniqueness: { scope: :provider }
  validates :unconfirmed_email, format: { with: URI::MailTo::EMAIL_REGEXP, allow_blank: true }

  def confirm!
    return false unless unconfirmed_or_reconfirming?

    return false if unconfirmed_email.present? && !update(email: unconfirmed_email, unconfirmed_email: nil)

    update(confirmed_at: Time.current)
  end

  def confirmable_email
    unconfirmed_email.presence || email
  end

  def reconfirming?
    unconfirmed_email.present?
  end

  def unconfirmed_or_reconfirming?
    unconfirmed? || reconfirming?
  end

  def confirmed?
    confirmed_at.present?
  end

  def generate_confirmation_token
    signed_id expires_in: CONFIRMATION_TOKEN_EXPIRATION, purpose: :confirm_email
  end

  def unconfirmed?
    !confirmed?
  end

  def send_confirmation_email!
    confirmation_token = generate_confirmation_token
    UserMailer.confirmation(self, confirmation_token).deliver_now
  end

  def generate_password_reset_token
    signed_id expires_in: PASSWORD_RESET_TOKEN_EXPIRATION, purpose: :reset_password
  end

  def send_password_reset_email!
    password_reset_token = generate_password_reset_token
    UserMailer.password_reset(self, password_reset_token).deliver_now
  end

  private

  def generate_api_key
    self.key = SecureRandom.hex(32)
  end

  def set_provider
    self.provider = EMAIL_PASSWORD_PROVIDER
  end

  def downcase_email
    self.email = email.downcase
  end

  def downcase_unconfirmed_email
    return if unconfirmed_email.nil?

    self.unconfirmed_email = unconfirmed_email.downcase
  end
end

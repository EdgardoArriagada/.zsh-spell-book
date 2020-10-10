# rbenv ubuntu https://github.com/rbenv/rbenv
if [ -d "$HOME/.rbenv" ]; then

  __${zsb}_prepare_lazy_load "$ZSB_DIR/src/configurations/ruby-config/rbenv.init" \
    rbenv ruby gem irb bundle rspec

  __${zsb}_prepare_lazy_load "$ZSB_DIR/src/configurations/ruby-config/rails.init" \
    rails rake
fi


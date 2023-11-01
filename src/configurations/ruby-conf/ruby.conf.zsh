# rbenv ubuntu https://github.com/rbenv/rbenv
if [[ -d $HOME/.rbenv ]]; then
  ZSB_RUBY_CONF_DIR=$ZSB_DIR/src/configurations/ruby-conf

  ${zsb}.prepareLazyLoad \
    $ZSB_RUBY_CONF_DIR/rbenv.init \
    rbenv ruby gem irb bundle rspec yard solargraph

  ${zsb}.prepareLazyLoad \
    $ZSB_RUBY_CONF_DIR/rails.init \
    rails rake
fi


# rbenv ubuntu https://github.com/rbenv/rbenv
if [[ -d ${HOME}/.rbenv ]]; then
  ${zsb}.prepareLazyLoad ${ZSB_DIR}/src/configurations/ruby-conf/rbenv.init \
    rbenv ruby gem irb bundle rspec yard solargraph

  ${zsb}.prepareLazyLoad ${ZSB_DIR}/src/configurations/ruby-conf/rails.init \
    rails rake
fi


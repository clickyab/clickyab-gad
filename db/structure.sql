CREATE TABLE ads
(
    ad_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ad_size TINYINT(1) DEFAULT '0',
    u_id INT(11) DEFAULT '0',
    ad_name TEXT,
    ad_url TEXT,
    ad_code TEXT,
    ad_title TEXT,
    ad_body TEXT,
    ad_img VARCHAR(255),
    ad_status TINYINT(1) DEFAULT '0',
    ad_reject_reason VARCHAR(50),
    ad_ctr FLOAT DEFAULT '0.1',
    ad_conv MEDIUMINT(6) DEFAULT '0',
    ad_time INT(11) DEFAULT '0',
    ad_type TINYINT(1) DEFAULT '0',
    ad_mainText VARCHAR(128),
    ad_defineText VARCHAR(128),
    ad_textColor VARCHAR(10),
    ad_target VARCHAR(30),
    ad_attribute TEXT,
    ad_hash_attribute VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE INDEX ad_hash_attribute ON ads (ad_hash_attribute);
CREATE INDEX ad_size ON ads (ad_size, ad_status);
CREATE INDEX u_id ON ads (u_id);
CREATE TABLE ads_correpted
(
    ad_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ad_size TINYINT(1) DEFAULT '0',
    u_id INT(11) DEFAULT '0',
    ad_name TEXT,
    ad_url TEXT,
    ad_code TEXT,
    ad_title TEXT,
    ad_body TEXT,
    ad_img VARCHAR(255),
    ad_status TINYINT(1) DEFAULT '0',
    ad_reject_reason VARCHAR(50),
    ad_ctr FLOAT DEFAULT '0.1',
    ad_conv MEDIUMINT(6) DEFAULT '0',
    ad_time INT(11) DEFAULT '0',
    ad_type TINYINT(1) DEFAULT '0',
    ad_mainText VARCHAR(128),
    ad_defineText VARCHAR(128),
    ad_textColor VARCHAR(10),
    ad_target VARCHAR(30),
    ad_attribute TEXT,
    ad_hash_attribute VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE INDEX ad_hash_attribute ON ads_correpted (ad_hash_attribute);
CREATE INDEX ad_size ON ads_correpted (ad_size, ad_status);
CREATE INDEX u_id ON ads_correpted (u_id);
CREATE TABLE ads_frequency
(
    af_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ad_id INT(11) DEFAULT '0',
    cp_id INT(11) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    af_count_total MEDIUMINT(6) DEFAULT '0',
    af_page_event INT(11) DEFAULT '0',
    af_count_today MEDIUMINT(6) DEFAULT '0',
    af_date INT(8) DEFAULT '0'
);
CREATE INDEX ad_id ON ads_frequency (ad_id, cop_id);
CREATE INDEX cop_id ON ads_frequency (cop_id);
CREATE INDEX cop_id_2 ON ads_frequency (cop_id, af_count_today);
CREATE TABLE api_users
(
    api_users_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    u_id INT(11),
    api_users_password VARCHAR(128) COMMENT 'MD5',
    api_users_token VARCHAR(256),
    api_users_token_expire DATETIME,
    api_users_access_table TEXT NOT NULL
);
CREATE INDEX u_id ON api_users (u_id);
CREATE TABLE apps
(
    app_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    u_id INT(11) DEFAULT '0',
    app_token VARCHAR(200),
    app_name VARCHAR(200),
    en_app_name VARCHAR(255),
    app_package VARCHAR(200),
    am_id TINYINT(1) DEFAULT '0' COMMENT 'Market Id',
    app_minbid INT(11) DEFAULT '500',
    app_floor_cpm SMALLINT(4) DEFAULT '0',
    app_div FLOAT DEFAULT '1.7',
    app_status TINYINT(1) DEFAULT '0',
    app_review TINYINT(1) DEFAULT '0' COMMENT '0 => pending,1 => review,2 => repending',
    app_today_ctr INT(11) DEFAULT '0',
    app_today_imps INT(11) DEFAULT '0',
    app_today_clicks INT(11) DEFAULT '0',
    app_date INT(11) DEFAULT '0' ,
    app_cat VARCHAR(255),
    app_notapprovedreason VARCHAR(255),
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE INDEX app_token ON apps (app_token);
CREATE INDEX u_id ON apps (u_id);
CREATE TABLE apps_android_ver
(
    aav_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    aav_version INT(11) DEFAULT '0'
);
CREATE UNIQUE INDEX aav_android_version ON apps_android_ver (aav_version);
CREATE TABLE apps_bak
(
    app_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    u_id INT(11) DEFAULT '0',
    app_token VARCHAR(200),
    app_name VARCHAR(200),
    en_app_name VARCHAR(255) NOT NULL,
    app_package VARCHAR(200),
    am_id TINYINT(1) DEFAULT '0' COMMENT 'Market Id',
    app_minbid INT(11) DEFAULT '500',
    app_floor_cpm SMALLINT(4) DEFAULT '0',
    app_div FLOAT DEFAULT '1.6',
    app_status TINYINT(1) DEFAULT '0',
    app_today_ctr INT(11) DEFAULT '0',
    app_today_imps INT(11) DEFAULT '0',
    app_today_clicks INT(11) DEFAULT '0',
    app_date INT(11) DEFAULT '0',
    app_cat VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE INDEX app_token ON apps_bak (app_token);
CREATE INDEX u_id ON apps_bak (u_id);
CREATE TABLE apps_brand_models
(
    abm_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ab_id INT(11) DEFAULT '0',
    abm_model VARCHAR(255)
);
CREATE INDEX abm_model ON apps_brand_models (abm_model);
CREATE INDEX ab_id ON apps_brand_models (ab_id);
CREATE TABLE apps_brands
(
    ab_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ab_brand VARCHAR(100),
    ab_show TINYINT(4) DEFAULT '1',
    ab_count INT(12) DEFAULT '0'
);
CREATE UNIQUE INDEX ab_brand ON apps_brands (ab_brand);
CREATE TABLE apps_carriers
(
    ac_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ac_carrier VARCHAR(100),
    ac_show TINYINT(4) DEFAULT '1',
    ac_count INT(12) DEFAULT '0'
);
CREATE UNIQUE INDEX ac_carrier ON apps_carriers (ac_carrier);
CREATE TABLE apps_install
(
    api_id INT(11) NOT NULL,
    u_id INT(11) DEFAULT '0',
    api_token VARCHAR(200),
    api_name VARCHAR(200),
    api_package VARCHAR(200),
    api_status VARCHAR(200)
);
CREATE TABLE apps_langs
(
    al_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    al_lang VARCHAR(125),
    al_show TINYINT(4) DEFAULT '1',
    al_count INT(12)
);
CREATE UNIQUE INDEX al_lang ON apps_langs (al_lang);
CREATE TABLE apps_market
(
    am_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    am_market VARCHAR(100),
    am_market_persian VARCHAR(50),
    am_market_os VARCHAR(100)
);
CREATE TABLE apps_networks
(
    an_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    an_network VARCHAR(100),
    an_show TINYINT(4) DEFAULT '1',
    an_count INT(12) DEFAULT '0'
);
CREATE UNIQUE INDEX an_network ON apps_networks (an_network);
CREATE TABLE apps_potential
(
    send TINYINT(1) DEFAULT '0',
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    email VARCHAR(255),
    unsub TINYINT(1) DEFAULT '0'
);
CREATE TABLE audit_log
(
    id INT(10) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    role_id INT(11),
    user_id INT(11),
    impersonator INT(11),
    for_who INT(11),
    action CHAR(30),
    target_id INT(10) unsigned,
    target_type VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE INDEX audit_log_action_index ON audit_log (action);
CREATE INDEX audit_log_for_who_foreign ON audit_log (for_who);
CREATE INDEX audit_log_impersonator_foreign ON audit_log (impersonator);
CREATE INDEX audit_log_role_id_foreign ON audit_log (role_id);
CREATE INDEX audit_log_target_id_target_type_index ON audit_log (target_id, target_type);
CREATE INDEX audit_log_user_id_foreign ON audit_log (user_id);
CREATE TABLE audit_log_details
(
    id INT(10) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    audit_id INT(10) unsigned NOT NULL,
    data TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE INDEX audit_log_details_audit_id_foreign ON audit_log_details (audit_id);
CREATE TABLE billing
(
    bi_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    u_id INT(11) DEFAULT '0',
    income_id INT(11) DEFAULT '0',
    bi_is_crm TINYINT(1) DEFAULT '0',
    bi_title VARCHAR(255),
    bi_amount INT(11) DEFAULT '0',
    bi_type INT(11) DEFAULT '0',
    bi_balance INT(11) DEFAULT '0',
    bi_time INT(11) DEFAULT '0',
    bi_date INT(11) DEFAULT '0',
    bi_reason TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE INDEX u_id ON billing (u_id);
CREATE UNIQUE INDEX u_id_2 ON billing (u_id, income_id, bi_amount, bi_time);
CREATE TABLE campaigns
(
    cp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_type TINYINT(2) DEFAULT '0',
    cp_billing_type VARCHAR(4),
    u_id INT(11) DEFAULT '0',
    cp_name VARCHAR(255),
    cp_network TINYINT(1) DEFAULT '0',
    cp_placement VARCHAR(2550),
    cp_wfilter VARCHAR(2550),
    cp_retargeting VARCHAR(255),
    cp_frequency INT(3) DEFAULT '2',
    cp_segment_id INT(11) DEFAULT '0',
    cp_app_brand VARCHAR(200),
    cp_net_provider VARCHAR(200),
    cp_app_lang VARCHAR(200),
    cp_app_market INT(11),
    cp_web_mobile TINYINT(1) DEFAULT '0',
    cp_web TINYINT(1) DEFAULT '0',
    cp_application TINYINT(1) DEFAULT '0',
    cp_video TINYINT(1) DEFAULT '0',
    cp_apps_carriers VARCHAR(200),
    cp_longmap VARCHAR(200),
    cp_latmap VARCHAR(200),
    cp_radius INT(11) DEFAULT '0',
    cp_opt_ctr TINYINT(1) DEFAULT '0',
    cp_opt_conv TINYINT(1) DEFAULT '0',
    cp_opt_br TINYINT(1) DEFAULT '0',
    cp_gender TINYINT(1) DEFAULT '0',
    cp_alexa TINYINT(1) DEFAULT '0',
    cp_fatfinger TINYINT(1) DEFAULT '1',
    cp_under TINYINT(1) DEFAULT '0',
    cp_geos VARCHAR(200),
    cp_region VARCHAR(200),
    cp_country VARCHAR(200),
    cp_hoods VARCHAR(200),
    cp_isp_blacklist VARCHAR(200),
    cp_cat VARCHAR(200),
    cp_like_app VARCHAR(200),
    cp_app VARCHAR(2550),
    cp_app_filter VARCHAR(2550),
    cp_keywords TEXT,
    cp_platforms VARCHAR(100),
    cp_platform_version VARCHAR(200),
    cp_maxbid INT(11) DEFAULT '0',
    cp_weekly_budget INT(11) DEFAULT '0',
    cp_daily_budget INT(11) DEFAULT '0',
    cp_total_budget INT(11) DEFAULT '0',
    cp_weekly_spend INT(11) DEFAULT '0',
    cp_total_spend INT(11) DEFAULT '0',
    cp_today_spend INT(11) DEFAULT '0',
    cp_clicks INT(11) DEFAULT '0',
    cp_ctr FLOAT DEFAULT '0',
    cp_imps INT(11) DEFAULT '0',
    cp_cpm INT(11) DEFAULT '0',
    cp_cpa INT(11) DEFAULT '0',
    cp_cpc INT(11) DEFAULT '0',
    cp_conv INT(11) DEFAULT '0',
    cp_conv_rate FLOAT DEFAULT '0',
    cp_revenue INT(11) DEFAULT '0',
    cp_roi INT(4) DEFAULT '0',
    cp_start INT(11) DEFAULT '0',
    cp_end INT(11) DEFAULT '0',
    cp_status INT(11) DEFAULT '1',
    cp_lastupdate INT(11) DEFAULT '0',
    cp_hour_start TINYINT(4) DEFAULT '0',
    cp_hour_end TINYINT(4) DEFAULT '24',
    is_crm TINYINT(4) DEFAULT '0',
    cp_lock INT(11) DEFAULT '0' NOT NULL COMMENT 'determine if the campaign was created through crm',
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX cp_lock ON campaigns (cp_lock);
CREATE INDEX u_id_idx ON campaigns (u_id);
CREATE TABLE campaigns_ads
(
    ca_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    ad_id INT(11) DEFAULT '0',
    ca_status TINYINT(4) DEFAULT '1',
    ca_imps INT(11) DEFAULT '0',
    ca_cpm INT(11) DEFAULT '0',
    ca_cpc INT(11) DEFAULT '0',
    ca_clicks INT(11) DEFAULT '0',
    ca_ctr FLOAT DEFAULT '0.1',
    ca_conv TINYINT(4) DEFAULT '0',
    ca_conv_rate FLOAT DEFAULT '0',
    ca_cpa INT(11) DEFAULT '0',
    ca_spend INT(11) DEFAULT '0',
    ca_lastupdate INT(11) DEFAULT '0'
);
CREATE INDEX ad_id ON campaigns_ads (ad_id);
CREATE INDEX cp_id ON campaigns_ads (cp_id);
CREATE INDEX cp_id_2 ON campaigns_ads (cp_id, ca_status);
CREATE TABLE campaigns_bak
(
    cp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_type TINYINT(2) DEFAULT '0',
    cp_billing_type VARCHAR(4),
    u_id INT(11) DEFAULT '0',
    cp_name VARCHAR(255),
    cp_network TINYINT(1) DEFAULT '0',
    cp_placement VARCHAR(2550),
    cp_wfilter VARCHAR(2550),
    cp_retargeting VARCHAR(255),
    cp_frequency INT(3) DEFAULT '2',
    cp_segment_id INT(11) DEFAULT '0',
    cp_app_brand VARCHAR(200),
    cp_net_provider VARCHAR(200),
    cp_app_lang VARCHAR(200),
    cp_app_market INT(11),
    cp_web_mobile TINYINT(1) DEFAULT '0',
    cp_web TINYINT(1) DEFAULT '0',
    cp_application TINYINT(1) DEFAULT '0',
    cp_video TINYINT(1) DEFAULT '0',
    cp_apps_carriers VARCHAR(200),
    cp_longmap VARCHAR(200),
    cp_latmap VARCHAR(200),
    cp_radius INT(11) DEFAULT '0',
    cp_opt_ctr TINYINT(1) DEFAULT '0',
    cp_opt_conv TINYINT(1) DEFAULT '0',
    cp_opt_br TINYINT(1) DEFAULT '0',
    cp_gender TINYINT(1) DEFAULT '0',
    cp_alexa TINYINT(1) DEFAULT '0',
    cp_fatfinger TINYINT(1) DEFAULT '1',
    cp_under TINYINT(1) DEFAULT '0',
    cp_geos VARCHAR(200),
    cp_region VARCHAR(200),
    cp_country VARCHAR(200),
    cp_hoods VARCHAR(200),
    cp_isp_blacklist VARCHAR(200),
    cp_cat VARCHAR(200),
    cp_like_app VARCHAR(200),
    cp_app VARCHAR(2550),
    cp_app_filter VARCHAR(2550),
    cp_keywords TEXT,
    cp_platforms VARCHAR(100),
    cp_platform_version VARCHAR(200),
    cp_maxbid INT(11) DEFAULT '0',
    cp_weekly_budget INT(11) DEFAULT '0',
    cp_daily_budget INT(11) DEFAULT '0',
    cp_total_budget INT(11) DEFAULT '0',
    cp_weekly_spend INT(11) DEFAULT '0',
    cp_total_spend INT(11) DEFAULT '0',
    cp_today_spend INT(11) DEFAULT '0',
    cp_clicks INT(11) DEFAULT '0',
    cp_ctr FLOAT DEFAULT '0',
    cp_imps INT(11) DEFAULT '0',
    cp_cpm INT(11) DEFAULT '0',
    cp_cpa INT(11) DEFAULT '0',
    cp_cpc INT(11) DEFAULT '0',
    cp_conv INT(11) DEFAULT '0',
    cp_conv_rate FLOAT DEFAULT '0',
    cp_revenue INT(11) DEFAULT '0',
    cp_roi INT(4) DEFAULT '0',
    cp_start INT(11) DEFAULT '0',
    cp_end INT(11) DEFAULT '0',
    cp_status INT(11) DEFAULT '1',
    cp_lastupdate INT(11) DEFAULT '0',
    cp_hour_start TINYINT(4) DEFAULT '0',
    cp_hour_end TINYINT(4) DEFAULT '24',
    is_crm TINYINT(4) DEFAULT '0',
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE campaigns_interests
(
    cpin_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    in_id INT(11) DEFAULT '0'
);
CREATE INDEX cp_id ON campaigns_interests (cp_id);
CREATE TABLE campaigns_keywords
(
    cpk_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    k_id INT(11) DEFAULT '0'
);
CREATE INDEX cp_id ON campaigns_keywords (cp_id);
CREATE TABLE campaigns_locations
(
    cpl_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    location_id INT(11) DEFAULT '0'
);
CREATE INDEX cp_id ON campaigns_locations (cp_id);
CREATE TABLE campaigns_new
(
    cp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_type TINYINT(2) DEFAULT '0',
    cp_billing_type VARCHAR(4) DEFAULT 'cpm',
    u_id INT(11) DEFAULT '0',
    cp_name LONGTEXT,
    cp_network TINYINT(1) DEFAULT '0',
    cp_placement VARCHAR(2550),
    cp_wfilter VARCHAR(2550) DEFAULT '#',
    cp_retargeting VARCHAR(255),
    cp_frequency INT(3) DEFAULT '2',
    cp_segment_id INT(11) DEFAULT '0',
    cp_opt_ctr TINYINT(1) DEFAULT '0',
    cp_opt_conv TINYINT(1) DEFAULT '0',
    cp_opt_br TINYINT(1) DEFAULT '0',
    cp_gender TINYINT(1) DEFAULT '0',
    cp_alexa TINYINT(1) DEFAULT '0',
    cp_fatfinger TINYINT(1) DEFAULT '1',
    cp_under TINYINT(1) DEFAULT '0',
    cp_geos VARCHAR(200),
    cp_region VARCHAR(200),
    cp_hoods VARCHAR(200),
    cp_isp_blacklist VARCHAR(200) DEFAULT '#',
    cp_cat VARCHAR(200),
    cp_like_app VARCHAR(200),
    cp_app VARCHAR(2550),
    cp_app_filter VARCHAR(2550),
    cp_keywords TEXT,
    cp_platforms VARCHAR(100),
    cp_platform_version VARCHAR(200) COMMENT 'aav_id from   `apps_android_ver`  table',
    cp_maxbid INT(11) DEFAULT '0',
    cp_weekly_budget INT(11) DEFAULT '0',
    cp_daily_budget INT(11) DEFAULT '0',
    cp_total_budget INT(11) DEFAULT '0',
    cp_weekly_spend INT(11) DEFAULT '0',
    cp_total_spend INT(11) DEFAULT '0',
    cp_today_spend INT(11) DEFAULT '0',
    cp_clicks INT(11) DEFAULT '0',
    cp_ctr FLOAT DEFAULT '0',
    cp_imps INT(11) DEFAULT '0',
    cp_cpm INT(11) DEFAULT '0',
    cp_cpa INT(11) DEFAULT '0',
    cp_cpc INT(11) DEFAULT '0',
    cp_conversions INT(11) DEFAULT '0',
    cp_revenue INT(11) DEFAULT '0',
    cp_roi INT(4) DEFAULT '0',
    cp_start INT(11) DEFAULT '0',
    cp_end INT(11) DEFAULT '0',
    cp_status INT(11) DEFAULT '1',
    cp_lastupdate INT(11) DEFAULT '0',
    cp_hour_start TINYINT(4) DEFAULT '0' NOT NULL,
    cp_hour_end TINYINT(4) DEFAULT '24' NOT NULL,
    cp_app_brand VARCHAR(200) COMMENT 'ab_id from apps_brands table',
    cp_net_provider VARCHAR(200) COMMENT 'an_id from apps_networks',
    cp_app_lang VARCHAR(200) COMMENT 'al_id from `apps_langs` ',
    cp_app_market INT(11),
    cp_web_mobile TINYINT(1) NOT NULL,
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    cp_apps_carriers VARCHAR(200),
    cp_web TINYINT(1) NOT NULL,
    cp_application TINYINT(1) NOT NULL,
    longmap VARCHAR(200),
    latmap VARCHAR(200),
    radius INT(11)
);
CREATE TABLE campaigns_placement
(
    cpp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    cpp_status INT(11) DEFAULT '0'
);
CREATE INDEX cp_id ON campaigns_placement (cp_id);
CREATE INDEX w_id ON campaigns_placement (w_id);
CREATE TABLE campaigns_platform
(
    cpp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    platform_id INT(11) DEFAULT '0'
);
CREATE INDEX cp_id ON campaigns_platform (cp_id);
CREATE TABLE campaigns_retargeting
(
    cpr_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0'
);
CREATE INDEX cp_id ON campaigns_retargeting (cp_id);
CREATE INDEX w_id ON campaigns_retargeting (w_id);
CREATE TABLE campaigns_segments
(
    cs_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    seg_id INT(11) DEFAULT '0',
    cs_conversions INT(11) DEFAULT '0',
    cs_revenue INT(11) DEFAULT '0',
    cs_lastupdate INT(11) DEFAULT '0'
);
CREATE INDEX cp_id ON campaigns_segments (cp_id);
CREATE TABLE categories
(
    cat_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cat_title VARCHAR(255),
    cat_title_persian VARCHAR(255),
    cat_count_w INT(11) DEFAULT '0' NOT NULL,
    cat_count_a INT(11) DEFAULT '0' NOT NULL,
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE TABLE categories_old
(
    cat_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cat_code VARCHAR(10),
    cat_title VARCHAR(255),
    cat_parent VARCHAR(10),
    cat_title_persian VARCHAR(255),
    cat_active INT(11) NOT NULL,
    cat_count_w INT(5) DEFAULT '0' NOT NULL,
    cat_count_a INT(5) DEFAULT '0' NOT NULL
);
CREATE TABLE clicks
(
    c_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    c_winnerbid INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    app_id INT(11) DEFAULT '0',
    wp_id INT(11) DEFAULT '0',
    cp_id INT(11) DEFAULT '0',
    ca_id INT(11) DEFAULT '0',
    slot_id INT(11) DEFAULT '0',
    sla_id INT(11) DEFAULT '0',
    ad_id INT(11) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    imp_id INT(11) DEFAULT '0',
    c_status TINYINT(2) DEFAULT '0',
    c_ip VARCHAR(20),
    c_referaddress VARCHAR(255),
    c_parenturl VARCHAR(255),
    c_fast INT(12) DEFAULT '0',
    c_os TINYINT(1) DEFAULT '0',
    c_time INT(11) DEFAULT '0',
    c_date INT(11) DEFAULT '0'
);
CREATE INDEX app_id ON clicks (app_id, c_date);
CREATE INDEX ca_id ON clicks (ca_id, c_status, c_date);
CREATE INDEX c_date ON clicks (c_date);
CREATE INDEX sla_id ON clicks (sla_id, c_status, c_date);
CREATE INDEX w_id ON clicks (w_id, c_date);
CREATE TABLE clicks_conv
(
    cc_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    c_id INT(11) DEFAULT '0',
    c_winnerbid INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    app_id INT(11) DEFAULT '0',
    wp_id INT(11) DEFAULT '0',
    cp_id INT(11) DEFAULT '0',
    ca_id INT(11) DEFAULT '0',
    slot_id INT(11) DEFAULT '0',
    sla_id INT(11) DEFAULT '0',
    ad_id INT(11) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    imp_id INT(11) DEFAULT '0',
    c_status TINYINT(2) DEFAULT '0',
    c_ip VARCHAR(20),
    c_referaddress VARCHAR(255),
    c_parenturl VARCHAR(255),
    c_ua TEXT,
    c_fast INT(12) DEFAULT '0',
    c_os TINYINT(1) DEFAULT '0',
    c_time INT(11) DEFAULT '0',
    c_date INT(11) DEFAULT '0',
    c_action VARCHAR(255)
);
CREATE INDEX app_id ON clicks_conv (app_id);
CREATE INDEX c_date ON clicks_conv (c_date);
CREATE UNIQUE INDEX c_id ON clicks_conv (c_id);
CREATE INDEX sla_id_single ON clicks_conv (sla_id);
CREATE INDEX slot_id ON clicks_conv (slot_id);
CREATE INDEX w_id ON clicks_conv (w_id);
CREATE TABLE conversions
(
    conv_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cs_id INT(11) DEFAULT '0',
    seg_id INT(11) DEFAULT '0',
    seg_convvalue INT(11) DEFAULT '0' NOT NULL,
    cp_id INT(11) DEFAULT '0',
    ad_id INT(11) DEFAULT '0',
    ca_id INT(11) DEFAULT '0',
    imp_id INT(11) DEFAULT '0',
    c_id INT(11) DEFAULT '0',
    wp_id INT(11) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    conv_time INT(11) DEFAULT '0',
    conv_date INT(11) DEFAULT '0'
);
CREATE TABLE cookie_profiles
(
    cop_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cop_key VARCHAR(20),
    cop_email VARCHAR(100),
    cop_last_ip VARCHAR(50),
    cop_gender TINYINT(1) DEFAULT '0',
    cop_alexa TINYINT(1) DEFAULT '0',
    cop_os TINYINT(1) DEFAULT '0',
    cop_browser TINYINT(2) DEFAULT '0',
    cop_city SMALLINT(4) DEFAULT '0',
    cop_age TINYINT(1) DEFAULT '0',
    cop_keywords TEXT,
    cop_active_date INT(9) DEFAULT '0'
);
CREATE UNIQUE INDEX cop_key ON cookie_profiles (cop_key);
CREATE TABLE cookie_webpages
(
    cwp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    wp_id INT(11) DEFAULT '0',
    w_id INT(12) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    cwp_time INT(11) DEFAULT '0',
    cwp_date INT(8) DEFAULT '0'
);
CREATE INDEX cop_id ON cookie_webpages (cop_id);
CREATE INDEX cwp_date ON cookie_webpages (cwp_date);
CREATE INDEX wp_id ON cookie_webpages (wp_id);
CREATE TABLE cookie_websites
(
    cw_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    w_id INT(12) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    cw_time INT(11) DEFAULT '0',
    cw_date INT(11) DEFAULT '0'
);
CREATE INDEX cop_id ON cookie_websites (cop_id);
CREATE INDEX cw_date ON cookie_websites (cw_date);
CREATE INDEX w_id ON cookie_websites (w_id, cop_id);
CREATE TABLE country
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    iso CHAR(2) NOT NULL,
    name VARCHAR(80) NOT NULL,
    nicename VARCHAR(80) NOT NULL,
    iso3 CHAR(3),
    numcode SMALLINT(6),
    phonecode INT(5) NOT NULL
);
CREATE TABLE coupons
(
    cpn_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cpn_code VARCHAR(16) NOT NULL,
    cpn_value INT(11) NOT NULL,
    u_id INT(11) DEFAULT '0',
    cpn_date_used INT(11) DEFAULT '0' NOT NULL,
    cpn_date_expire INT(11) NOT NULL
);
CREATE TABLE cp_zero
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    q TEXT NOT NULL
);
CREATE TABLE crm_emails
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL
);
CREATE TABLE docker
(
    docker_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    docker_ip_client VARCHAR(255),
    docker_ip_server VARCHAR(255),
    docker_time INT(11) DEFAULT '0'
);
CREATE TABLE finder_admin
(
    id INT(10) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    pass CHAR(32) NOT NULL
);
CREATE TABLE finder_cells
(
    id INT(10) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cellname VARCHAR(6) NOT NULL,
    top_left_lat VARCHAR(50) NOT NULL,
    top_left_long VARCHAR(50) NOT NULL,
    bottom_left_lat VARCHAR(50) NOT NULL,
    bottom_left_long VARCHAR(50) NOT NULL,
    bottom_right_lat VARCHAR(50) NOT NULL,
    bottom_right_long VARCHAR(50) NOT NULL,
    top_right_lat VARCHAR(50) NOT NULL,
    top_right_long VARCHAR(50) NOT NULL
);
CREATE INDEX theindex ON finder_cells (top_left_lat, bottom_left_lat, top_left_long, bottom_right_long);
CREATE TABLE finder_city_parts
(
    id MEDIUMINT(7) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    title VARCHAR(4) NOT NULL,
    cellgroup TEXT NOT NULL
);
CREATE TABLE finder_logs
(
    id INT(20) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT(10) unsigned DEFAULT '0',
    cell_id INT(10) unsigned DEFAULT '0',
    imei VARCHAR(50),
    android_id VARCHAR(50),
    carrier VARCHAR(35),
    mcc MEDIUMINT(6) unsigned DEFAULT '0',
    mnc MEDIUMINT(6) unsigned DEFAULT '0',
    lac INT(10) unsigned DEFAULT '0',
    cid INT(10) unsigned DEFAULT '0',
    ip VARCHAR(16),
    l_time INT(10) unsigned DEFAULT '0',
    locations VARCHAR(101),
    time INT(10) unsigned
);
CREATE TABLE finder_logs_sdk_old
(
    id INT(20) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cell_id INT(10) unsigned DEFAULT '0',
    android_id VARCHAR(50),
    android_version VARCHAR(20),
    parameters TEXT,
    carrier VARCHAR(35),
    mcc MEDIUMINT(6) unsigned DEFAULT '0',
    mnc MEDIUMINT(6) unsigned DEFAULT '0',
    lac INT(10) unsigned DEFAULT '0',
    cid INT(10) unsigned DEFAULT '0',
    ip VARCHAR(16),
    locations VARCHAR(101),
    time INT(10) unsigned
);
CREATE INDEX mcc ON finder_logs_sdk_old (mcc, mnc, lac, cid);
CREATE TABLE finder_logs_sdk_true
(
    id INT(20) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cell_id INT(10) DEFAULT '0',
    recheck TINYINT(1) DEFAULT '0',
    android_id VARCHAR(50),
    android_version VARCHAR(20),
    parameters TEXT,
    carrier VARCHAR(35),
    mcc MEDIUMINT(6) unsigned DEFAULT '0',
    mnc MEDIUMINT(6) unsigned DEFAULT '0',
    lac INT(10) unsigned DEFAULT '0',
    cid INT(10) unsigned DEFAULT '0',
    ip VARCHAR(100),
    locations VARCHAR(255),
    time INT(10) unsigned
);
CREATE INDEX cell_id ON finder_logs_sdk_true (cell_id, locations);
CREATE INDEX mcc ON finder_logs_sdk_true (mcc, mnc, lac, cid);
CREATE TABLE finder_users
(
    id INT(10) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    android_id VARCHAR(50) NOT NULL,
    imei VARCHAR(50) NOT NULL,
    carrier VARCHAR(50) NOT NULL,
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    time INT(10) unsigned NOT NULL,
    ip VARCHAR(16) NOT NULL,
    status TINYINT(2) unsigned DEFAULT '1' NOT NULL
);
CREATE TABLE google_users
(
    google_id DECIMAL(21) PRIMARY KEY NOT NULL,
    google_name VARCHAR(60) NOT NULL,
    google_email VARCHAR(60) NOT NULL,
    google_link VARCHAR(60) NOT NULL,
    google_picture_link VARCHAR(60) NOT NULL
);
CREATE TABLE hits
(
    hit_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    w_id INT(11) DEFAULT '0',
    hit_date INT(11) DEFAULT '0'
);
CREATE INDEX w_id ON hits (w_id);
CREATE TABLE impressions
(
    imp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    wp_id INT(11) DEFAULT '0',
    app_id INT(11) DEFAULT '0',
    ad_id INT(11) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    ca_id INT(11) DEFAULT '0',
    imp_ipaddress VARCHAR(50),
    imp_referaddress TEXT,
    imp_parenturl TEXT,
    imp_url TEXT,
    imp_winnerbid INT(11) DEFAULT '0',
    imp_status TINYINT(1) DEFAULT '0',
    imp_cookie TINYINT(1) DEFAULT '1',
    imp_alexa TINYINT(1) DEFAULT '0',
    imp_conv TINYINT(1) DEFAULT '0',
    imp_flash TINYINT(1) DEFAULT '1',
    imp_time INT(11) DEFAULT '0',
    imp_date INT(8) DEFAULT '0'
);
CREATE INDEX app_id ON impressions (app_id, imp_date);
CREATE INDEX ca_id ON impressions (ca_id, imp_date);
CREATE INDEX imp_date ON impressions (imp_date);
CREATE INDEX w_id ON impressions (w_id, imp_date);
CREATE TABLE `impressions-cells`
(
    imp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    app_id INT(11) DEFAULT '0',
    wp_id INT(11) DEFAULT '0',
    ad_id INT(11) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    ca_id INT(11) DEFAULT '0',
    slot_id INT(11) DEFAULT '0',
    sla_id INT(11) DEFAULT '0',
    cell_id INT(11) DEFAULT '0',
    imp_ipaddress VARCHAR(50),
    imp_referaddress TEXT,
    imp_parenturl TEXT,
    imp_url TEXT,
    imp_winnerbid INT(11) DEFAULT '0',
    imp_status TINYINT(1) DEFAULT '0',
    imp_cookie TINYINT(1) DEFAULT '1',
    imp_alexa TINYINT(1) DEFAULT '0',
    imp_flash TINYINT(1) DEFAULT '1',
    imp_time INT(11) DEFAULT '0',
    imp_date INT(8) DEFAULT '0'
);
CREATE INDEX app_id ON `impressions-cells` (app_id);
CREATE INDEX ca_id ON `impressions-cells` (ca_id);
CREATE INDEX imp_date ON `impressions-cells` (imp_date);
CREATE INDEX sla_id ON `impressions-cells` (sla_id);
CREATE INDEX slot_id ON `impressions-cells` (slot_id);
CREATE INDEX w_id ON `impressions-cells` (w_id);
CREATE TABLE interests
(
    in_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    in_parent_id INT(11) DEFAULT '0',
    in_gender TINYINT(1) DEFAULT '0',
    in_age TINYINT(1) DEFAULT '0',
    in_name INT(11)
);
CREATE TABLE invoices
(
    in_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    in_accept TINYINT(2) DEFAULT '0' NOT NULL,
    in_serial INT(11) DEFAULT '0' NOT NULL,
    in_date INT(11) DEFAULT '0' NOT NULL,
    u_id INT(11) NOT NULL,
    in_price INT(11) NOT NULL,
    in_title VARCHAR(200),
    in_user_register_id INT(11) NOT NULL,
    in_sale_condition INT(1) DEFAULT '0'
);
CREATE TABLE invoices_details
(
    ind_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    in_id INT(11) NOT NULL,
    ind_description VARCHAR(500),
    ind_timing VARCHAR(255),
    ind_price INT(11) NOT NULL,
    ind_price_off INT(11) NOT NULL
);
CREATE TABLE ip2location
(
    ip_from INT(10) unsigned,
    ip_to INT(10) unsigned,
    country_code CHAR(2),
    country_name VARCHAR(64),
    region_name VARCHAR(128),
    city_name VARCHAR(128)
);
CREATE INDEX CC ON ip2location (country_code, ip_to);
CREATE TABLE ip2location_old
(
    ip_from INT(10) unsigned,
    ip_to INT(10) unsigned,
    country_code CHAR(2),
    country_name VARCHAR(64),
    region_name VARCHAR(128),
    city_name VARCHAR(128)
);
CREATE INDEX idx_ip_from ON ip2location_old (ip_from);
CREATE INDEX idx_ip_from_to ON ip2location_old (ip_from, ip_to);
CREATE INDEX idx_ip_to ON ip2location_old (ip_to);
CREATE TABLE ip_to_country_remove
(
    IP_FROM BIGINT(20) unsigned NOT NULL,
    IP_TO BIGINT(20) unsigned NOT NULL,
    REGISTRY CHAR(7) NOT NULL,
    ASSIGNED BIGINT(20) NOT NULL,
    CTRY CHAR(2) NOT NULL,
    CNTRY CHAR(3) NOT NULL,
    COUNTRY VARCHAR(100) NOT NULL,
    CONSTRAINT `PRIMARY` PRIMARY KEY (IP_FROM, IP_TO)
);
CREATE TABLE jabe_abzar_imps
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    imp_id INT(11) NOT NULL,
    date INT(11) NOT NULL
);
CREATE TABLE keywords
(
    k_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    k_string VARCHAR(255),
    k_string_md5 VARCHAR(64),
    k_confirm TINYINT(1) DEFAULT '0',
    k_count INT(11) DEFAULT '0'
);
CREATE UNIQUE INDEX k_string_md5 ON keywords (k_string_md5);
CREATE TABLE keywords_interests
(
    ki_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    k_id INT(11) DEFAULT '0',
    in_id INT(11) DEFAULT '0'
);
CREATE INDEX in_id ON keywords_interests (in_id);
CREATE INDEX k_id ON keywords_interests (k_id);
CREATE TABLE keywords_webpages
(
    kwp_id INT(11) PRIMARY KEY NOT NULL,
    k_id INT(11) DEFAULT '0',
    wp_id INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0'
);
CREATE INDEX k_id ON keywords_webpages (k_id, wp_id);
CREATE INDEX k_id_2 ON keywords_webpages (k_id);
CREATE INDEX wp_id ON keywords_webpages (wp_id);
CREATE TABLE list_browser
(
    browser_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    browser_value VARCHAR(100),
    browser_name VARCHAR(100)
);
CREATE UNIQUE INDEX bor_value ON list_browser (browser_value);
CREATE TABLE list_city
(
    location_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    location_name TEXT,
    location_name_persian TEXT,
    location_master MEDIUMINT(6) DEFAULT '0',
    location_select TINYINT(1) DEFAULT '0',
    location_code INT(11) NOT NULL,
    location_country VARCHAR(3),
    location_region INT(11) DEFAULT '0' NOT NULL
);
CREATE TABLE list_locations
(
    location_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    location_name VARCHAR (100),
    location_name_persian TEXT,
    location_master MEDIUMINT(6) DEFAULT '0',
    location_select TINYINT(1) DEFAULT '0'
);
CREATE INDEX location_name ON list_locations (location_name);
CREATE TABLE list_platform
(
    platform_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    platform_network TINYINT(1) DEFAULT '0',
    platform_value VARCHAR(100),
    platform_name VARCHAR(100)
);
CREATE UNIQUE INDEX osl_value ON list_platform (platform_value);
CREATE TABLE list_region
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    region_name VARCHAR(50),
    region_name_persian VARCHAR(100),
    region_code INT(2) DEFAULT '0'
);
CREATE TABLE list_region2
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    region_name VARCHAR(50),
    region_name_persian VARCHAR(100),
    region_code INT(2) DEFAULT '0'
);
CREATE TABLE neighborhoods
(
    id MEDIUMINT(7) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    title TEXT,
    cellsgroup TEXT
);
CREATE TABLE neighborhoods_old
(
    id MEDIUMINT(7) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    cellsgroup TEXT NOT NULL
);
CREATE TABLE password_resets
(
    u_email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX password_resets_email_index ON password_resets (u_email);
CREATE INDEX password_resets_token_index ON password_resets (token);
CREATE TABLE payment_transaction
(
    pt_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    u_id INT(11) DEFAULT '0',
    pt_amount INT(11) DEFAULT '0',
    pt_type TINYINT(2) DEFAULT '0',
    pt_gate VARCHAR(100),
    pt_status TINYINT(2) DEFAULT '0',
    pt_authority VARCHAR(255),
    pt_refid VARCHAR(255),
    pt_time INT(11) DEFAULT '0',
    pt_date INT(11) DEFAULT '0',
    pt_flag TINYINT(2) DEFAULT '0'
);
CREATE INDEX pt_authority ON payment_transaction (pt_authority);
CREATE INDEX u_id ON payment_transaction (u_id);
CREATE TABLE permission_role
(
    permission_id INT(10) unsigned NOT NULL,
    role_id INT(10) unsigned NOT NULL,
    CONSTRAINT `PRIMARY` PRIMARY KEY (permission_id, role_id)
);
CREATE INDEX permission_role_role_id_foreign ON permission_role (role_id);
CREATE TABLE permissions
(
    id INT(10) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    label VARCHAR(255),
    access VARCHAR(100) DEFAULT 'list' NOT NULL,
    action VARCHAR(100) DEFAULT 'own' NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE TABLE qlog
(
    q_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    q_content TEXT,
    q_time INT(11) DEFAULT '0'
);
CREATE TABLE role_user
(
    role_id INT(10) unsigned NOT NULL,
    user_id INT(10) unsigned NOT NULL,
    CONSTRAINT `PRIMARY` PRIMARY KEY (role_id, user_id)
);
CREATE INDEX role_user_user_id_foreign ON role_user (user_id);
CREATE TABLE roles
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    label VARCHAR(255),
    childes VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE TABLE segments
(
    seg_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    u_id INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    seg_pattern VARCHAR(255),
    seg_type TINYINT(1) DEFAULT '0',
    seg_name VARCHAR(255),
    seg_isconv TINYINT(1) DEFAULT '0',
    seg_convvalue INT(11) DEFAULT '0',
    seg_conversions INT(11) DEFAULT '0',
    seg_visitors INT(11) DEFAULT '0',
    seg_pageview INT(11) DEFAULT '0',
    seg_lastupdate INT(11) DEFAULT '0'
);
CREATE INDEX u_id ON segments (u_id);
CREATE TABLE slots
(
    slot_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    slot_pubilc_id BIGINT(20) DEFAULT '0',
    slot_name VARCHAR(255),
    slot_size TINYINT(2) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    app_id INT(11) DEFAULT '0',
    slot_avg_daily_imps INT(11) DEFAULT '0',
    slot_avg_daily_clicks INT(11) DEFAULT '0',
    slot_floor_cpm INT(11) DEFAULT '0',
    slot_total_monthly_cost INT(11) DEFAULT '0',
    slot_lastupdate INT(11) DEFAULT '0',
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE INDEX slot_pubilc_id ON slots (slot_pubilc_id, w_id);
CREATE TABLE slots_ads
(
    sla_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    slot_id INT(11) DEFAULT '0',
    ad_id INT(11) DEFAULT '0',
    sla_imps INT(11) DEFAULT '0',
    sla_clicks INT(11) DEFAULT '0',
    sla_ctr FLOAT DEFAULT '0',
    sla_conv INT(11) DEFAULT '0',
    sla_conv_rate FLOAT DEFAULT '0',
    sla_cpa INT(11) DEFAULT '0',
    sla_cpm INT(11) DEFAULT '0',
    sla_spend INT(11) DEFAULT '0',
    sla_lastupdate INT(11) DEFAULT '0'
);
CREATE INDEX slot_id ON slots_ads (slot_id, ad_id, sla_ctr);
CREATE TABLE slots_bak
(
    slot_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    slot_pubilc_id BIGINT(20) DEFAULT '0',
    slot_name VARCHAR(255),
    slot_size TINYINT(2) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    app_id INT(11) DEFAULT '0',
    slot_avg_daily_imps INT(11) DEFAULT '0',
    slot_avg_daily_clicks INT(11) DEFAULT '0',
    slot_floor_cpm INT(11) DEFAULT '0',
    slot_total_monthly_cost INT(11) DEFAULT '0',
    slot_lastupdate INT(11) DEFAULT '0',
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE INDEX slot_pubilc_id ON slots_bak (slot_pubilc_id, w_id);
CREATE TABLE statistics_ads
(
    sa_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ad_id INT(11) DEFAULT '0',
    cp_id INT(11) DEFAULT '0',
    ca_id INT(11) DEFAULT '0',
    sa_clicks INT(11) DEFAULT '0',
    sa_imps INT(11) DEFAULT '0',
    sa_ctr FLOAT DEFAULT '0',
    sa_conv INT(11) DEFAULT '0',
    sa_conv_rate FLOAT DEFAULT '0',
    sa_cpa INT(11) DEFAULT '0',
    sa_spend INT(11) DEFAULT '0',
    sa_day INT(11) DEFAULT '0',
    sa_done TINYINT(1) DEFAULT '0'
);
CREATE INDEX ca_id ON statistics_ads (ca_id, sa_day);
CREATE INDEX cp_id ON statistics_ads (ca_id);
CREATE INDEX sa_day ON statistics_ads (sa_day);
CREATE TABLE statistics_alexa
(
    sal_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    sal_ip VARCHAR(100),
    sal_ref VARCHAR(255),
    sal_user_agent VARCHAR(255),
    sal_show TINYINT(1) DEFAULT '0',
    sal_time INT(11) DEFAULT '0',
    sal_day INT(11) DEFAULT '0'
);
CREATE TABLE statistics_apps
(
    sa_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    app_id INT(11) DEFAULT '0',
    sa_clicks INT(11) DEFAULT '0',
    sa_allclicks INT(11) DEFAULT '0',
    sa_imps INT(11) DEFAULT '0',
    sa_ctr FLOAT DEFAULT '0',
    sa_conv INT(11) DEFAULT '0',
    sa_conv_rate FLOAT DEFAULT '0',
    sa_cpa INT(11) DEFAULT '0',
    sa_spend INT(11) DEFAULT '0',
    sa_day INT(11) DEFAULT '0',
    sa_done TINYINT(4) DEFAULT '0'
);
CREATE UNIQUE INDEX app_id ON statistics_apps (app_id, sa_day);
CREATE INDEX w_id ON statistics_apps (app_id);
CREATE TABLE statistics_campaigns
(
    sc_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    sc_clicks INT(11) DEFAULT '0',
    sc_imps INT(11) DEFAULT '0',
    sc_ctr FLOAT DEFAULT '0',
    sc_conv INT(11) DEFAULT '0',
    sc_conv_rate FLOAT DEFAULT '0',
    sc_cpa INT(11) DEFAULT '0',
    sc_spend INT(11) DEFAULT '0',
    sc_day INT(11) DEFAULT '0',
    sc_done TINYINT(4) DEFAULT '0'
);
CREATE INDEX cp_id ON statistics_campaigns (cp_id);
CREATE TABLE statistics_segments
(
    ss_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    seg_id INT(11) DEFAULT '0',
    cs_id INT(11) DEFAULT '0',
    ss_pageviews INT(11) DEFAULT '0',
    ss_newvisitors INT(11) DEFAULT '0',
    ss_conversions INT(11) DEFAULT '0',
    ss_revenue INT(11) DEFAULT '0',
    ss_day INT(11) DEFAULT '0',
    ss_done TINYINT(4) DEFAULT '0'
);
CREATE INDEX seg_id ON statistics_segments (seg_id);
CREATE TABLE statistics_slot
(
    ssl_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ad_id INT(11) DEFAULT '0',
    slot_id INT(11) DEFAULT '0',
    sla_id INT(11) DEFAULT '0',
    ssl_clicks INT(11) DEFAULT '0',
    ssl_imps INT(11) DEFAULT '0',
    ssl_ctr FLOAT DEFAULT '0',
    ssl_conv INT(11) DEFAULT '0',
    ssl_conv_rate FLOAT DEFAULT '0',
    ssl_cpa INT(11) DEFAULT '0',
    ssl_spend INT(11) DEFAULT '0',
    ssl_day INT(11) DEFAULT '0',
    ssl_done TINYINT(1) DEFAULT '0'
);
CREATE INDEX sla_id ON statistics_slot (sla_id, ssl_day);
CREATE INDEX slot_id ON statistics_slot (slot_id);
CREATE INDEX slot_id_day ON statistics_slot (slot_id, ssl_day);
CREATE INDEX ssl_day ON statistics_slot (ssl_day);
CREATE TABLE statistics_websites
(
    sw_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    w_id INT(11) DEFAULT '0',
    sw_clicks INT(11) DEFAULT '0',
    sw_allclicks INT(11) DEFAULT '0',
    sw_imps INT(11) DEFAULT '0',
    sw_ctr FLOAT DEFAULT '0',
    sw_conv INT(11) DEFAULT '0',
    sw_conv_rate FLOAT DEFAULT '0',
    sw_cpa INT(11) DEFAULT '0',
    sw_spend INT(11) DEFAULT '0',
    sw_day INT(11) DEFAULT '0',
    sw_done TINYINT(4) DEFAULT '0'
);
CREATE INDEX w_id ON statistics_websites (w_id);
CREATE UNIQUE INDEX w_id_2 ON statistics_websites (w_id, sw_day);
CREATE TABLE tickets
(
    ti_id INT(10) unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ti_slug VARCHAR(255),
    u_id INT(11) NOT NULL,
    ti_title VARCHAR(255),
    ti_content LONGTEXT NOT NULL,
    ti_response LONGTEXT NOT NULL,
    ti_status TINYINT(4) DEFAULT '0' NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE TABLE trueview
(
    tv_click_id INT(11) PRIMARY KEY NOT NULL
);
CREATE TABLE unsubscribe
(
    email VARCHAR(255)
);
CREATE TABLE users
(
    u_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    u_email VARCHAR(100),
    u_password VARCHAR(255),
    u_access_key VARCHAR(200),
    u_profile_type TINYINT(1) DEFAULT '0',
    u_role TINYINT(1) DEFAULT '0' COMMENT 'general = 0 , 1 = super admin , 2 => publisher, 3 => accounting , 4 => sells , 5 => support',
    u_firstname VARCHAR(64),
    u_lastname VARCHAR(64),
    u_fullname VARCHAR(255),
    u_avatar VARCHAR(200),
    u_mobile VARCHAR(20),
    u_phone VARCHAR(20),
    u_postcode VARCHAR(255),
    u_address TEXT,
    u_company_name VARCHAR(255),
    u_melli_code VARCHAR(255),
    u_account_number VARCHAR(255),
    u_card_number VARCHAR(255),
    u_sheba_number VARCHAR(255),
    u_account_holder VARCHAR(255),
    u_bank_name VARCHAR(255),
    u_province VARCHAR(200),
    u_city VARCHAR(200),
    u_economic_number VARCHAR(200),
    u_register_number VARCHAR(200),
    u_mobile_confirm TINYINT(1) DEFAULT '0',
    u_email_confirm TINYINT(1) DEFAULT '0',
    u_mobile_confirm_string VARCHAR(12),
    u_email_confirm_string VARCHAR(12),
    u_balance INT(11) DEFAULT '0',
    u_today_spend INT(11) DEFAULT '0',
    u_close TINYINT(1) DEFAULT '0',
    u_unsubscribe TINYINT(1) DEFAULT '0',
    u_nlsend TINYINT(1) DEFAULT '0',
    u_ref VARCHAR(200),
    u_ip VARCHAR(100),
    u_date DATE,
    u_time INT(11) DEFAULT '0',
    u_google_id VARCHAR(25),
    u_gender VARCHAR(4),
    u_longitude VARCHAR(10),
    u_latitude VARCHAR(10),
    u_owner VARCHAR(128),
    u_originatinglead VARCHAR(128),
    u_customer_code INT(12),
    u_gid VARCHAR(64),
    u_leadstatus INT(1) DEFAULT '0' NOT NULL,
    u_parent INT(11) DEFAULT '0',
    u_crm TINYINT(1) DEFAULT '0',
    remember_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    bReadByCRM TINYINT(1) DEFAULT '0' NOT NULL
);
CREATE INDEX u_access_key ON users (u_access_key);
CREATE UNIQUE INDEX u_email ON users (u_email);
CREATE INDEX u_email_2 ON users (u_email, u_password);
CREATE TABLE users_bak2
(
    u_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    u_email VARCHAR(100),
    u_password VARCHAR(255),
    u_access_key VARCHAR(200),
    u_profile_type TINYINT(1) DEFAULT '0',
    u_role TINYINT(1) DEFAULT '0' COMMENT 'general = 0 , 1 = super admin , 2 => publisher, 3 => accounting , 4 => sells , 5 => support',
    u_firstname VARCHAR(64),
    u_lastname VARCHAR(64),
    u_fullname VARCHAR(255),
    u_avatar VARCHAR(200),
    u_mobile VARCHAR(20),
    u_phone VARCHAR(20),
    u_postcode VARCHAR(255),
    u_address TEXT,
    u_company_name VARCHAR(255),
    u_melli_code VARCHAR(255),
    u_account_number VARCHAR(255),
    u_card_number VARCHAR(255),
    u_sheba_number VARCHAR(255),
    u_account_holder VARCHAR(255),
    u_bank_name VARCHAR(255),
    u_province VARCHAR(200),
    u_city VARCHAR(200),
    u_economic_number VARCHAR(200),
    u_register_number VARCHAR(200),
    u_mobile_confirm TINYINT(1) DEFAULT '0',
    u_email_confirm TINYINT(1) DEFAULT '0',
    u_mobile_confirm_string VARCHAR(12),
    u_email_confirm_string VARCHAR(12),
    u_balance INT(11) DEFAULT '0',
    u_today_spend INT(11) DEFAULT '0',
    u_close TINYINT(1) DEFAULT '0',
    u_unsubscribe TINYINT(1) DEFAULT '0',
    u_nlsend TINYINT(1) DEFAULT '0',
    u_ref VARCHAR(200),
    u_ip VARCHAR(100),
    u_date DATE,
    u_time INT(11) DEFAULT '0',
    u_google_id VARCHAR(25),
    u_gender VARCHAR(4),
    u_longitude VARCHAR(10),
    u_latitude VARCHAR(10),
    u_owner VARCHAR(128),
    u_originatinglead VARCHAR(128),
    u_customer_code INT(12),
    u_gid VARCHAR(64),
    u_leadstatus INT(1) DEFAULT '0' NOT NULL,
    u_parent INT(11) DEFAULT '0',
    u_crm TINYINT(1) DEFAULT '0',
    remember_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE INDEX u_access_key ON users_bak2 (u_access_key);
CREATE UNIQUE INDEX u_email ON users_bak2 (u_email);
CREATE INDEX u_email_2 ON users_bak2 (u_email, u_password);
CREATE TABLE users_log
(
    ul_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    ul_email VARCHAR(100) NOT NULL,
    ul_password VARCHAR(255),
    ul_access_key VARCHAR(200),
    ul_profile_type TINYINT(1) DEFAULT '0',
    ul_role TINYINT(1) DEFAULT '0' COMMENT 'general = 0 , 1 = super admin , 2 => publisher, 3 => accounting , 4 => sells , 5 => support',
    ul_fullname VARCHAR(255),
    ul_avatar VARCHAR(200),
    ul_mobile VARCHAR(20) DEFAULT '09',
    ul_phone VARCHAR(20),
    ul_postcode VARCHAR(255),
    ul_address TEXT,
    ul_company_name VARCHAR(255),
    ul_melli_code VARCHAR(255),
    ul_account_number VARCHAR(255),
    ul_card_number VARCHAR(255),
    ul_sheba_number VARCHAR(255) DEFAULT 'IR',
    ul_account_holder VARCHAR(255),
    ul_bank_name VARCHAR(255),
    ul_province VARCHAR(200),
    ul_city VARCHAR(200),
    ul_economic_number VARCHAR(200),
    ul_register_number VARCHAR(200),
    ul_mobile_confirm TINYINT(1) DEFAULT '0',
    ul_mobile_confirm_string VARCHAR(12) DEFAULT '0',
    ul_email_confirm_string VARCHAR(12) DEFAULT '0',
    ul_email_confirm TINYINT(1) DEFAULT '0',
    ul_balance INT(11) DEFAULT '0',
    ul_close TINYINT(1) DEFAULT '0',
    ul_unsubscribe TINYINT(1) DEFAULT '0',
    ul_nlsend TINYINT(1) DEFAULT '0',
    ul_ref VARCHAR(200),
    ul_ip VARCHAR(100),
    ul_date DATE,
    ul_time INT(11) DEFAULT '0',
    ul_google_id VARCHAR(25),
    ul_gender VARCHAR(4)
);
CREATE UNIQUE INDEX ul_email ON users_log (ul_email);
CREATE TABLE webpages
(
    wp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    w_id INT(11) DEFAULT '0',
    wp_url VARCHAR(255),
    wp_md5 VARCHAR(64),
    wp_keywords TEXT
);
CREATE INDEX wp_md5 ON webpages (wp_md5);
CREATE INDEX w_id ON webpages (w_id);
CREATE TABLE websites
(
    w_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    u_id INT(11) DEFAULT '0',
    w_pub_id BIGINT(16) DEFAULT '0',
    w_domain VARCHAR(100),
    w_name VARCHAR(100),
    w_categories VARCHAR(255),
    w_profile_type TINYINT(4) DEFAULT '0',
    w_minbid INT(11) DEFAULT '1500',
    w_floor_cpm INT(5) DEFAULT '2100',
    w_status TINYINT(4) DEFAULT '0' COMMENT '0 => pending,1 => accepted ,2 => rejected, 3 => deleted',
    w_review TINYINT(1) DEFAULT '0' COMMENT '0 => pending,1 => review,2 => repending',
    w_alexarank INT(11) DEFAULT '0',
    w_div FLOAT DEFAULT '1.7',
    w_mobad TINYINT(1) DEFAULT '0',
    w_nativead TINYINT(1) DEFAULT '0',
    w_fatfinger TINYINT(1) DEFAULT '0',
    w_publish_start INT(11) DEFAULT '0',
    w_publish_end INT(11) DEFAULT '0',
    w_publish_cost INT(11) DEFAULT '0',
    w_prepayment TINYINT(1) DEFAULT '0',
    w_today_ctr FLOAT DEFAULT '0',
    w_today_imps INT(11) DEFAULT '0',
    w_today_clicks INT(11) DEFAULT '0',
    w_date INT(11) DEFAULT '0',
    w_notapprovedreason VARCHAR(255),
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE INDEX u_id ON websites (u_id);
CREATE INDEX w_domain ON websites (w_domain, w_today_imps);
CREATE UNIQUE INDEX w_pub_id ON websites (w_pub_id);
CREATE INDEX w_pub_id_2 ON websites (w_status, w_pub_id);
CREATE TABLE websites_categories
(
    wc_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    w_id INT(11) DEFAULT '0',
    cat_id INT(11) DEFAULT '0'
);
CREATE INDEX w_id ON websites_categories (w_id);
CREATE TABLE websites_potential
(
    wpt_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    wpt_domain VARCHAR(200),
    wpt_fullname VARCHAR(128),
    wpt_status TINYINT(4),
    wpt_alexa INT(11) DEFAULT '0',
    wpt_email VARCHAR(255),
    wpt_phone VARCHAR(255),
    wpt_telephone VARCHAR(30),
    wpt_partner VARCHAR(255),
    wpt_categories VARCHAR(255),
    wpt_clickyab_id INT(11) DEFAULT '0',
    wpt_description TEXT,
    wpt_telegram VARCHAR(200),
    wpt_category VARCHAR(512),
    created_at TIMESTAMP DEFAULT '0000-00-00 00:00:00',
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE INDEX wpt_domain ON websites_potential (wpt_domain);
CREATE TABLE withdrawal
(
    wd_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    u_id INT(11) DEFAULT '0',
    wd_amount INT(11) DEFAULT '0',
    wd_card_number VARCHAR(255),
    wd_sheba_number VARCHAR(255),
    wd_tracking_code VARCHAR(255),
    wd_open TINYINT(1) DEFAULT '1',
    wd_date INT(11) DEFAULT '0',
    wd_date_paid INT(11) DEFAULT '0',
    wd_reason VARCHAR(255),
    wd_uid_approver INT(10),
    wd_date_approvement VARCHAR(100),
    wd_approve INT(1) DEFAULT '0' COMMENT '0 => Not Approved ; 1 => Approved ; ',
    wd_is_advertiser INT(1) DEFAULT '0' COMMENT '1 = > is Advertiser ; 0 => publisher ;',
    wd_not_approved_reasons VARCHAR(255),
    wd_review_date INT(11),
    wd_review_description VARCHAR(255),
    wd_can_request_review INT(1) DEFAULT '0',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00'
);
CREATE INDEX u_id ON withdrawal (u_id);
CREATE TABLE impressions20161108
(
    imp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    app_id INT(11) DEFAULT '0',
    wp_id INT(11) DEFAULT '0',
    ad_id INT(11) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    ca_id INT(11) DEFAULT '0',
    slot_id INT(11) DEFAULT '0',
    sla_id INT(11) DEFAULT '0',
    cell_id INT(11) DEFAULT '0',
    hood_id INT(11) DEFAULT '0',
    imp_ipaddress VARCHAR(50),
    imp_referaddress TEXT,
    imp_parenturl TEXT,
    imp_url TEXT,
    imp_winnerbid INT(11) DEFAULT '0',
    imp_status TINYINT(1) DEFAULT '0',
    imp_cookie TINYINT(1) DEFAULT '1',
    imp_alexa TINYINT(1) DEFAULT '0',
    imp_flash TINYINT(1) DEFAULT '1',
    imp_time INT(11) DEFAULT '0',
    imp_date INT(8) DEFAULT '0'
);
CREATE INDEX app_id ON impressions20161108 (app_id);
CREATE INDEX ca_id ON impressions20161108 (ca_id);
CREATE INDEX imp_date ON impressions20161108 (imp_date);
CREATE INDEX sla_id ON impressions20161108 (sla_id);
CREATE INDEX slot_id ON impressions20161108 (slot_id);
CREATE INDEX w_id ON impressions20161108 (w_id);
CREATE TABLE impressions20161109
(
    imp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    app_id INT(11) DEFAULT '0',
    wp_id INT(11) DEFAULT '0',
    ad_id INT(11) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    ca_id INT(11) DEFAULT '0',
    slot_id INT(11) DEFAULT '0',
    sla_id INT(11) DEFAULT '0',
    cell_id INT(11) DEFAULT '0',
    hood_id INT(11) DEFAULT '0',
    imp_ipaddress VARCHAR(50),
    imp_referaddress TEXT,
    imp_parenturl TEXT,
    imp_url TEXT,
    imp_winnerbid INT(11) DEFAULT '0',
    imp_status TINYINT(1) DEFAULT '0',
    imp_cookie TINYINT(1) DEFAULT '1',
    imp_alexa TINYINT(1) DEFAULT '0',
    imp_flash TINYINT(1) DEFAULT '1',
    imp_time INT(11) DEFAULT '0',
    imp_date INT(8) DEFAULT '0'
);
CREATE INDEX app_id ON impressions20161109 (app_id);
CREATE INDEX ca_id ON impressions20161109 (ca_id);
CREATE INDEX imp_date ON impressions20161109 (imp_date);
CREATE INDEX sla_id ON impressions20161109 (sla_id);
CREATE INDEX slot_id ON impressions20161109 (slot_id);
CREATE INDEX w_id ON impressions20161109 (w_id);
CREATE TABLE impressions20161110
(
    imp_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cp_id INT(11) DEFAULT '0',
    w_id INT(11) DEFAULT '0',
    app_id INT(11) DEFAULT '0',
    wp_id INT(11) DEFAULT '0',
    ad_id INT(11) DEFAULT '0',
    cop_id INT(11) DEFAULT '0',
    ca_id INT(11) DEFAULT '0',
    slot_id INT(11) DEFAULT '0',
    sla_id INT(11) DEFAULT '0',
    cell_id INT(11) DEFAULT '0',
    hood_id INT(11) DEFAULT '0',
    imp_ipaddress VARCHAR(50),
    imp_referaddress TEXT,
    imp_parenturl TEXT,
    imp_url TEXT,
    imp_winnerbid INT(11) DEFAULT '0',
    imp_status TINYINT(1) DEFAULT '0',
    imp_cookie TINYINT(1) DEFAULT '1',
    imp_alexa TINYINT(1) DEFAULT '0',
    imp_flash TINYINT(1) DEFAULT '1',
    imp_time INT(11) DEFAULT '0',
    imp_date INT(8) DEFAULT '0'
);
CREATE INDEX app_id ON impressions20161110 (app_id);
CREATE INDEX ca_id ON impressions20161110 (ca_id);
CREATE INDEX imp_date ON impressions20161110 (imp_date);
CREATE INDEX sla_id ON impressions20161110 (sla_id);
CREATE INDEX slot_id ON impressions20161110 (slot_id);
CREATE INDEX w_id ON impressions20161110 (w_id);

alter table slots_ads  add unique index index_name (slot_id,ad_id);
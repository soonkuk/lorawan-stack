// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttipb

import fmt "fmt"

func (dst *Configuration) SetFields(src *Configuration, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "default_cluster":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster
				if (src == nil || src.DefaultCluster == nil) && dst.DefaultCluster == nil {
					continue
				}
				if src != nil {
					newSrc = src.DefaultCluster
				}
				if dst.DefaultCluster != nil {
					newDst = dst.DefaultCluster
				} else {
					newDst = &Configuration_Cluster{}
					dst.DefaultCluster = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DefaultCluster = src.DefaultCluster
				} else {
					dst.DefaultCluster = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_UI) SetFields(src *Configuration_UI, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "branding_base_url":
			if len(subs) > 0 {
				return fmt.Errorf("'branding_base_url' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.BrandingBaseURL = src.BrandingBaseURL
			} else {
				var zero string
				dst.BrandingBaseURL = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster) SetFields(src *Configuration_Cluster, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "ui":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_UI
				if (src == nil || src.UI == nil) && dst.UI == nil {
					continue
				}
				if src != nil {
					newSrc = src.UI
				}
				if dst.UI != nil {
					newDst = dst.UI
				} else {
					newDst = &Configuration_UI{}
					dst.UI = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.UI = src.UI
				} else {
					dst.UI = nil
				}
			}
		case "is":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer
				if (src == nil || src.IS == nil) && dst.IS == nil {
					continue
				}
				if src != nil {
					newSrc = src.IS
				}
				if dst.IS != nil {
					newDst = dst.IS
				} else {
					newDst = &Configuration_Cluster_IdentityServer{}
					dst.IS = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.IS = src.IS
				} else {
					dst.IS = nil
				}
			}
		case "ns":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_NetworkServer
				if (src == nil || src.NS == nil) && dst.NS == nil {
					continue
				}
				if src != nil {
					newSrc = src.NS
				}
				if dst.NS != nil {
					newDst = dst.NS
				} else {
					newDst = &Configuration_Cluster_NetworkServer{}
					dst.NS = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.NS = src.NS
				} else {
					dst.NS = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer) SetFields(src *Configuration_Cluster_IdentityServer, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "user_registration":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_UserRegistration
				if (src == nil || src.UserRegistration == nil) && dst.UserRegistration == nil {
					continue
				}
				if src != nil {
					newSrc = src.UserRegistration
				}
				if dst.UserRegistration != nil {
					newDst = dst.UserRegistration
				} else {
					newDst = &Configuration_Cluster_IdentityServer_UserRegistration{}
					dst.UserRegistration = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.UserRegistration = src.UserRegistration
				} else {
					dst.UserRegistration = nil
				}
			}
		case "profile_picture":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_ProfilePicture
				if (src == nil || src.ProfilePicture == nil) && dst.ProfilePicture == nil {
					continue
				}
				if src != nil {
					newSrc = src.ProfilePicture
				}
				if dst.ProfilePicture != nil {
					newDst = dst.ProfilePicture
				} else {
					newDst = &Configuration_Cluster_IdentityServer_ProfilePicture{}
					dst.ProfilePicture = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ProfilePicture = src.ProfilePicture
				} else {
					dst.ProfilePicture = nil
				}
			}
		case "end_device_picture":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_EndDevicePicture
				if (src == nil || src.EndDevicePicture == nil) && dst.EndDevicePicture == nil {
					continue
				}
				if src != nil {
					newSrc = src.EndDevicePicture
				}
				if dst.EndDevicePicture != nil {
					newDst = dst.EndDevicePicture
				} else {
					newDst = &Configuration_Cluster_IdentityServer_EndDevicePicture{}
					dst.EndDevicePicture = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.EndDevicePicture = src.EndDevicePicture
				} else {
					dst.EndDevicePicture = nil
				}
			}
		case "user_rights":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_UserRights
				if (src == nil || src.UserRights == nil) && dst.UserRights == nil {
					continue
				}
				if src != nil {
					newSrc = src.UserRights
				}
				if dst.UserRights != nil {
					newDst = dst.UserRights
				} else {
					newDst = &Configuration_Cluster_IdentityServer_UserRights{}
					dst.UserRights = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.UserRights = src.UserRights
				} else {
					dst.UserRights = nil
				}
			}
		case "oauth":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_OAuth
				if (src == nil || src.OAuth == nil) && dst.OAuth == nil {
					continue
				}
				if src != nil {
					newSrc = src.OAuth
				}
				if dst.OAuth != nil {
					newDst = dst.OAuth
				} else {
					newDst = &Configuration_Cluster_IdentityServer_OAuth{}
					dst.OAuth = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OAuth = src.OAuth
				} else {
					dst.OAuth = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_NetworkServer) SetFields(src *Configuration_Cluster_NetworkServer, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "dev_addr_prefixes":
			if len(subs) > 0 {
				return fmt.Errorf("'dev_addr_prefixes' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DevAddrPrefixes = src.DevAddrPrefixes
			} else {
				dst.DevAddrPrefixes = nil
			}
		case "deduplication_window":
			if len(subs) > 0 {
				return fmt.Errorf("'deduplication_window' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DeduplicationWindow = src.DeduplicationWindow
			} else {
				dst.DeduplicationWindow = nil
			}
		case "cooldown_window":
			if len(subs) > 0 {
				return fmt.Errorf("'cooldown_window' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CooldownWindow = src.CooldownWindow
			} else {
				dst.CooldownWindow = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_UserRegistration) SetFields(src *Configuration_Cluster_IdentityServer_UserRegistration, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "invitation":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_UserRegistration_Invitation
				if (src == nil || src.Invitation == nil) && dst.Invitation == nil {
					continue
				}
				if src != nil {
					newSrc = src.Invitation
				}
				if dst.Invitation != nil {
					newDst = dst.Invitation
				} else {
					newDst = &Configuration_Cluster_IdentityServer_UserRegistration_Invitation{}
					dst.Invitation = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Invitation = src.Invitation
				} else {
					dst.Invitation = nil
				}
			}
		case "contact_info_validation":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_UserRegistration_ContactInfoValidation
				if (src == nil || src.ContactInfoValidation == nil) && dst.ContactInfoValidation == nil {
					continue
				}
				if src != nil {
					newSrc = src.ContactInfoValidation
				}
				if dst.ContactInfoValidation != nil {
					newDst = dst.ContactInfoValidation
				} else {
					newDst = &Configuration_Cluster_IdentityServer_UserRegistration_ContactInfoValidation{}
					dst.ContactInfoValidation = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ContactInfoValidation = src.ContactInfoValidation
				} else {
					dst.ContactInfoValidation = nil
				}
			}
		case "admin_approval":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_UserRegistration_AdminApproval
				if (src == nil || src.AdminApproval == nil) && dst.AdminApproval == nil {
					continue
				}
				if src != nil {
					newSrc = src.AdminApproval
				}
				if dst.AdminApproval != nil {
					newDst = dst.AdminApproval
				} else {
					newDst = &Configuration_Cluster_IdentityServer_UserRegistration_AdminApproval{}
					dst.AdminApproval = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.AdminApproval = src.AdminApproval
				} else {
					dst.AdminApproval = nil
				}
			}
		case "password_requirements":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_UserRegistration_PasswordRequirements
				if (src == nil || src.PasswordRequirements == nil) && dst.PasswordRequirements == nil {
					continue
				}
				if src != nil {
					newSrc = src.PasswordRequirements
				}
				if dst.PasswordRequirements != nil {
					newDst = dst.PasswordRequirements
				} else {
					newDst = &Configuration_Cluster_IdentityServer_UserRegistration_PasswordRequirements{}
					dst.PasswordRequirements = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.PasswordRequirements = src.PasswordRequirements
				} else {
					dst.PasswordRequirements = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_ProfilePicture) SetFields(src *Configuration_Cluster_IdentityServer_ProfilePicture, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "disable_upload":
			if len(subs) > 0 {
				return fmt.Errorf("'disable_upload' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DisableUpload = src.DisableUpload
			} else {
				dst.DisableUpload = nil
			}
		case "use_gravatar":
			if len(subs) > 0 {
				return fmt.Errorf("'use_gravatar' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UseGravatar = src.UseGravatar
			} else {
				dst.UseGravatar = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_EndDevicePicture) SetFields(src *Configuration_Cluster_IdentityServer_EndDevicePicture, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "disable_upload":
			if len(subs) > 0 {
				return fmt.Errorf("'disable_upload' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DisableUpload = src.DisableUpload
			} else {
				dst.DisableUpload = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_UserRights) SetFields(src *Configuration_Cluster_IdentityServer_UserRights, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "create_applications":
			if len(subs) > 0 {
				return fmt.Errorf("'create_applications' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CreateApplications = src.CreateApplications
			} else {
				dst.CreateApplications = nil
			}
		case "create_clients":
			if len(subs) > 0 {
				return fmt.Errorf("'create_clients' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CreateClients = src.CreateClients
			} else {
				dst.CreateClients = nil
			}
		case "create_gateways":
			if len(subs) > 0 {
				return fmt.Errorf("'create_gateways' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CreateGateways = src.CreateGateways
			} else {
				dst.CreateGateways = nil
			}
		case "create_organizations":
			if len(subs) > 0 {
				return fmt.Errorf("'create_organizations' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CreateOrganizations = src.CreateOrganizations
			} else {
				dst.CreateOrganizations = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_OAuth) SetFields(src *Configuration_Cluster_IdentityServer_OAuth, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "providers":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_OAuth_AuthProviders
				if (src == nil || src.Providers == nil) && dst.Providers == nil {
					continue
				}
				if src != nil {
					newSrc = src.Providers
				}
				if dst.Providers != nil {
					newDst = dst.Providers
				} else {
					newDst = &Configuration_Cluster_IdentityServer_OAuth_AuthProviders{}
					dst.Providers = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Providers = src.Providers
				} else {
					dst.Providers = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_UserRegistration_Invitation) SetFields(src *Configuration_Cluster_IdentityServer_UserRegistration_Invitation, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "required":
			if len(subs) > 0 {
				return fmt.Errorf("'required' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Required = src.Required
			} else {
				dst.Required = nil
			}
		case "token_ttl":
			if len(subs) > 0 {
				return fmt.Errorf("'token_ttl' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TokenTTL = src.TokenTTL
			} else {
				dst.TokenTTL = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_UserRegistration_ContactInfoValidation) SetFields(src *Configuration_Cluster_IdentityServer_UserRegistration_ContactInfoValidation, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "required":
			if len(subs) > 0 {
				return fmt.Errorf("'required' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Required = src.Required
			} else {
				dst.Required = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_UserRegistration_AdminApproval) SetFields(src *Configuration_Cluster_IdentityServer_UserRegistration_AdminApproval, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "required":
			if len(subs) > 0 {
				return fmt.Errorf("'required' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Required = src.Required
			} else {
				dst.Required = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_UserRegistration_PasswordRequirements) SetFields(src *Configuration_Cluster_IdentityServer_UserRegistration_PasswordRequirements, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "min_length":
			if len(subs) > 0 {
				return fmt.Errorf("'min_length' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.MinLength = src.MinLength
			} else {
				dst.MinLength = nil
			}
		case "max_length":
			if len(subs) > 0 {
				return fmt.Errorf("'max_length' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.MaxLength = src.MaxLength
			} else {
				dst.MaxLength = nil
			}
		case "min_uppercase":
			if len(subs) > 0 {
				return fmt.Errorf("'min_uppercase' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.MinUppercase = src.MinUppercase
			} else {
				dst.MinUppercase = nil
			}
		case "min_digits":
			if len(subs) > 0 {
				return fmt.Errorf("'min_digits' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.MinDigits = src.MinDigits
			} else {
				dst.MinDigits = nil
			}
		case "min_special":
			if len(subs) > 0 {
				return fmt.Errorf("'min_special' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.MinSpecial = src.MinSpecial
			} else {
				dst.MinSpecial = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_OAuth_AuthProviders) SetFields(src *Configuration_Cluster_IdentityServer_OAuth_AuthProviders, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "oidc":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_OAuth_AuthProviders_OpenIDConnect
				if (src == nil || src.OIDC == nil) && dst.OIDC == nil {
					continue
				}
				if src != nil {
					newSrc = src.OIDC
				}
				if dst.OIDC != nil {
					newDst = dst.OIDC
				} else {
					newDst = &Configuration_Cluster_IdentityServer_OAuth_AuthProviders_OpenIDConnect{}
					dst.OIDC = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OIDC = src.OIDC
				} else {
					dst.OIDC = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_OAuth_AuthProviders_Shared) SetFields(src *Configuration_Cluster_IdentityServer_OAuth_AuthProviders_Shared, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "name":
			if len(subs) > 0 {
				return fmt.Errorf("'name' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Name = src.Name
			} else {
				var zero string
				dst.Name = zero
			}
		case "allow_registrations":
			if len(subs) > 0 {
				return fmt.Errorf("'allow_registrations' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AllowRegistrations = src.AllowRegistrations
			} else {
				var zero bool
				dst.AllowRegistrations = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Configuration_Cluster_IdentityServer_OAuth_AuthProviders_OpenIDConnect) SetFields(src *Configuration_Cluster_IdentityServer_OAuth_AuthProviders_OpenIDConnect, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "shared":
			if len(subs) > 0 {
				var newDst, newSrc *Configuration_Cluster_IdentityServer_OAuth_AuthProviders_Shared
				if src != nil {
					newSrc = &src.Configuration_Cluster_IdentityServer_OAuth_AuthProviders_Shared
				}
				newDst = &dst.Configuration_Cluster_IdentityServer_OAuth_AuthProviders_Shared
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Configuration_Cluster_IdentityServer_OAuth_AuthProviders_Shared = src.Configuration_Cluster_IdentityServer_OAuth_AuthProviders_Shared
				} else {
					var zero Configuration_Cluster_IdentityServer_OAuth_AuthProviders_Shared
					dst.Configuration_Cluster_IdentityServer_OAuth_AuthProviders_Shared = zero
				}
			}
		case "client_id":
			if len(subs) > 0 {
				return fmt.Errorf("'client_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ClientID = src.ClientID
			} else {
				var zero string
				dst.ClientID = zero
			}
		case "client_secret":
			if len(subs) > 0 {
				return fmt.Errorf("'client_secret' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ClientSecret = src.ClientSecret
			} else {
				var zero string
				dst.ClientSecret = zero
			}
		case "redirect_url":
			if len(subs) > 0 {
				return fmt.Errorf("'redirect_url' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.RedirectURL = src.RedirectURL
			} else {
				var zero string
				dst.RedirectURL = zero
			}
		case "provider_url":
			if len(subs) > 0 {
				return fmt.Errorf("'provider_url' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ProviderURL = src.ProviderURL
			} else {
				var zero string
				dst.ProviderURL = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

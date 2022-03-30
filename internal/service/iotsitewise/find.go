package iotsitewise

// func FindThingByName(conn *iot.IoT, name string) (*iot.DescribeThingOutput, error) {
// 	input := &iot.DescribeThingInput{
// 		ThingName: aws.String(name),
// 	}

// 	output, err := conn.DescribeThing(input)

// 	if tfawserr.ErrCodeEquals(err, iot.ErrCodeResourceNotFoundException) {
// 		return nil, &resource.NotFoundError{
// 			LastError:   err,
// 			LastRequest: input,
// 		}
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	if output == nil {
// 		return nil, tfresource.NewEmptyResultError(input)
// 	}

// 	return output, nil
// }

// func FindThingGroupByName(conn *iot.IoT, name string) (*iot.DescribeThingGroupOutput, error) {
// 	input := &iot.DescribeThingGroupInput{
// 		ThingGroupName: aws.String(name),
// 	}

// 	output, err := conn.DescribeThingGroup(input)

// 	if tfawserr.ErrCodeEquals(err, iot.ErrCodeResourceNotFoundException) {
// 		return nil, &resource.NotFoundError{
// 			LastError:   err,
// 			LastRequest: input,
// 		}
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	if output == nil {
// 		return nil, tfresource.NewEmptyResultError(input)
// 	}

// 	return output, nil
// }

// func FindThingGroupMembership(conn *iot.IoT, thingGroupName, thingName string) error {
// 	input := &iot.ListThingGroupsForThingInput{
// 		ThingName: aws.String(thingName),
// 	}

// 	var v *iot.GroupNameAndArn

// 	err := conn.ListThingGroupsForThingPages(input, func(page *iot.ListThingGroupsForThingOutput, lastPage bool) bool {
// 		if page == nil {
// 			return !lastPage
// 		}

// 		for _, group := range page.ThingGroups {
// 			if aws.StringValue(group.GroupName) == thingGroupName {
// 				v = group

// 				return false
// 			}
// 		}

// 		return !lastPage
// 	})

// 	if tfawserr.ErrCodeEquals(err, iot.ErrCodeResourceNotFoundException) {
// 		return &resource.NotFoundError{
// 			LastError:   err,
// 			LastRequest: input,
// 		}
// 	}

// 	if v == nil {
// 		return tfresource.NewEmptyResultError(input)
// 	}

// 	return nil
// }

// aws iotsitewise list-access-policies --resource-type "PORTAL" --resource-id 2e1b4315-94e1-433b-aaf9-0472261ec1db

// {
//     "accessPolicySummaries": [
//         {
//             "id": "34177956-a3e0-4279-aff8-8794c7db867e",
//             "identity": {
//                 "user": {
//                     "id": "c3673d5d4b-8995a246-9149-44dc-978e-262cc1f2d137"
//                 }
//             },
//             "resource": {
//                 "portal": {
//                     "id": "2e1b4315-94e1-433b-aaf9-0472261ec1db"
//                 }
//             },
//             "permission": "ADMINISTRATOR",
//             "creationDate": "2022-03-30T17:26:46+02:00",
//             "lastUpdateDate": "2022-03-30T17:26:46+02:00"
//         }
//     ]
// }

// aws iotsitewise delete-access-policy --access-policy-id 34177956-a3e0-4279-aff8-8794c7db867e